package main

import (
	"github.com/spf13/cobra"
	"rag-new/internal/migrations"
)

func init() {
	rootCmd.AddCommand(migrateUpCmd, migrateDownCmd)
}

var migrateUpCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := CreateApp()
		if err != nil {
			panic(err)
			return
		}
		migrations.NewMigrate(app.X)
		migrations.Migrate()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:  "migrate:rollback",
	Long: "rollback database",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := CreateApp()
		if err != nil {
			panic(err)
			return
		}
		migrations.NewMigrate(app.X)
		migrations.Rollback()
	},
}
