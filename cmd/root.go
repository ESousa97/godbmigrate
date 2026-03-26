package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lucassousa/godbmigrate/internal/db"
	"github.com/spf13/cobra"
)

var (
	dsn   string
	debug bool
	store *db.MigrationStore
)

var rootCmd = &cobra.Command{
	Use:   "godbmigrate",
	Short: "A simple and robust database migration tool",
	Long:  `godbmigrate is a CLI tool designed to manage database schema evolutions with ease and safety.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		if debug {
			opts.Level = slog.LevelDebug
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
		slog.SetDefault(logger)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "PostgreSQL DSN")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
}

// initDB initializes the database connection if a DSN is provided and the command needs it
func initDB() error {
	var err error
	store, err = db.Connect(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}
