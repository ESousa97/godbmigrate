package cmd

import (
	"fmt"
	"os"

	"github.com/lucassousa/godbmigrate/internal/db"
	"github.com/spf13/cobra"
)

var (
	dsn   string
	store *db.MigrationStore
)

var rootCmd = &cobra.Command{
	Use:   "godbmigrate",
	Short: "godbmigrate is a simple tool to manage database migrations",
	Long:  `A fast and flexible database migration tool for Go projects.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "PostgreSQL DSN")
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
