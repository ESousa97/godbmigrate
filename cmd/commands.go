package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const migrationsDir = "migrations"

var revertAll bool

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new pair of migration files",
	Long:  `Generates a new set of .up.sql and .down.sql migration files with a timestamped prefix.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")

		if err := os.MkdirAll(migrationsDir, 0755); err != nil {
			slog.Error("Failed to create migrations directory", "error", err)
			return
		}

		upFile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
		downFile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

		if err := createFile(upFile); err != nil {
			slog.Error("Failed to create up file", "error", err)
			return
		}
		if err := createFile(downFile); err != nil {
			slog.Error("Failed to create down file", "error", err)
			return
		}

		slog.Info("Migration files created", "up", upFile, "down", downFile)
	},
}

func createFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available migration files",
	Long:  `Scans the migrations directory and lists all discovered migration files in alphabetical order.`,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := os.ReadDir(migrationsDir)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Warn("No migrations folder found")
				return
			}
			slog.Error("Error reading migrations", "error", err)
			return
		}

		var fileNames []string
		for _, f := range files {
			if !f.IsDir() {
				fileNames = append(fileNames, f.Name())
			}
		}

		sort.Strings(fileNames)

		if len(fileNames) == 0 {
			slog.Info("No migrations found")
			return
		}

		fmt.Println("Migrations:")
		for _, name := range fileNames {
			fmt.Printf("- %s\n", name)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current migration version",
	Long:  `Queries the database to identify the last successfully applied migration version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initDB(); err != nil {
			return err
		}
		defer func() {
			_ = store.Close()
		}()

		version, err := store.GetLatestVersion()
		if err != nil {
			return err
		}

		if version == 0 {
			slog.Info("No migrations have been applied yet")
		} else {
			slog.Info("Current status", "version", version)
		}

		return nil
	},
}

type migrationFile struct {
	version int64
	name    string
	path    string
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Long:  `Identifies migration files that haven't been applied yet and executes them in chronological order.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initDB(); err != nil {
			return err
		}
		defer func() {
			_ = store.Close()
		}()

		if err := store.AcquireLock(); err != nil {
			return err
		}
		defer func() {
			_ = store.ReleaseLock()
		}()

		currentVersion, err := store.GetLatestVersion()
		if err != nil {
			return err
		}

		files, err := os.ReadDir(migrationsDir)
		if err != nil {
			return fmt.Errorf("could not read migrations directory: %w", err)
		}

		var pending []migrationFile
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".up.sql") {
				continue
			}

			parts := strings.Split(f.Name(), "_")
			if len(parts) < 2 {
				continue
			}

			version, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				continue
			}

			if version > currentVersion {
				pending = append(pending, migrationFile{
					version: version,
					name:    f.Name(),
					path:    filepath.Join(migrationsDir, f.Name()),
				})
			}
		}

		sort.Slice(pending, func(i, j int) bool {
			return pending[i].version < pending[j].version
		})

		if len(pending) == 0 {
			slog.Info("No pending migrations to apply")
			return nil
		}

		for _, m := range pending {
			slog.Info("Applying migration", "name", m.name)

			content, err := os.ReadFile(m.path)
			if err != nil {
				slog.Error("Could not read migration file", "name", m.name, "error", err)
				return err
			}

			if err := store.ApplyMigration(m.version, string(content)); err != nil {
				slog.Error("Migration failed", "name", m.name, "error", err)
				return err
			}

			slog.Info("Migration applied successfully", "name", m.name)
		}

		slog.Info("Migration process completed", "count", len(pending))
		return nil
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert migrations",
	Long:  `Reverts the last applied migration. If the --all flag is provided, it reverts all applied migrations in reverse order.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initDB(); err != nil {
			return err
		}
		defer func() {
			_ = store.Close()
		}()

		if err := store.AcquireLock(); err != nil {
			return err
		}
		defer func() {
			_ = store.ReleaseLock()
		}()

		appliedVersions, err := store.GetAppliedVersions()
		if err != nil {
			return err
		}

		if len(appliedVersions) == 0 {
			slog.Info("No applied migrations to revert")
			return nil
		}

		var toRevert []int64
		if revertAll {
			toRevert = appliedVersions
		} else {
			toRevert = []int64{appliedVersions[0]}
		}

		files, err := os.ReadDir(migrationsDir)
		if err != nil {
			return fmt.Errorf("could not read migrations directory: %w", err)
		}

		for _, version := range toRevert {
			downFile := ""
			for _, f := range files {
				if f.IsDir() || !strings.HasSuffix(f.Name(), ".down.sql") {
					continue
				}
				if strings.HasPrefix(f.Name(), strconv.FormatInt(version, 10)+"_") {
					downFile = f.Name()
					break
				}
			}

			if downFile == "" {
				slog.Error("Down migration file not found", "version", version)
				return fmt.Errorf("down migration file not found for version %d", version)
			}

			slog.Info("Reverting migration", "name", downFile)

			path := filepath.Join(migrationsDir, downFile)
			content, err := os.ReadFile(path)
			if err != nil {
				slog.Error("Could not read migration file", "name", downFile, "error", err)
				return err
			}

			if err := store.RevertMigration(version, string(content)); err != nil {
				slog.Error("Reversion failed", "name", downFile, "error", err)
				return err
			}

			slog.Info("Migration reverted successfully", "name", downFile)
		}

		slog.Info("Reversion process completed", "count", len(toRevert))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(upCmd)
	
	downCmd.Flags().BoolVar(&revertAll, "all", false, "Revert all applied migrations")
	rootCmd.AddCommand(downCmd)
}
