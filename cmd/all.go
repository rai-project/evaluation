package cmd

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use: "all",
	Aliases: []string{
		"eval_all",
	},
	Short: "Get all evaluation information from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFlags := cmd.Parent().Flags()
		originalDatabaseName := databaseName
		for _, pcmd := range allCmds {
			databaseName = originalDatabaseName
			pcmd.Flags().AddFlagSet(rootFlags)
			log.WithField("command", pcmd.Name()).
				WithField("args", args).
				Info("running evaluation command")
			pp.Println("running evaluation command ", args, "  ", pcmd.Name())
			pcmd.SilenceUsage = true
			err := pcmd.PreRunE(pcmd, args)
			if err != nil {
				log.WithError(err).
					WithField("command", pcmd.Name()).
					WithField("args", args).
					Error("failed to pre run evaluation command")
				continue
			}
			err = pcmd.RunE(pcmd, args)
			if err != nil {
				log.WithError(err).
					WithField("command", pcmd.Name()).
					WithField("args", args).
					Error("failed to run evaluation command")
				continue
			}
		}

		return nil
	},
}
