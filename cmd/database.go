package cmd

import "github.com/spf13/cobra"

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database operations",
}

func init() {
	databaseCmd.AddCommand(divergenceCmds...)
}
