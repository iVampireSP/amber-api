package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rag",
	Short: "rag",
	Long:  "rag",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			panic(err)
			return
		}
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
		return
	}
}
