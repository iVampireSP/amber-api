package main

import (
	"github.com/spf13/cobra"
	"rag-new/internal/migrations"
)

func init() {
	rootCmd.AddCommand(migrateUpCmd)
}

var migrateUpCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Long:  "Migrate database",
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
