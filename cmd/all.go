package cmd

import (
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use: "all",
	Aliases: []string{
		"eval_all",
	},
	Short: "Get all evaluation information from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, cmd := range allCmds {
			argWoFlags := cmd.Flags().Args()
			log.WithField("command", cmd.Name()).
				WithField("args", argWoFlags).
				Debug("running evaluation command")
			cmd.SilenceUsage = true
			err := cmd.PreRunE(cmd, argWoFlags)
			if err != nil {
				log.WithError(err).
					WithField("command", cmd.Name()).
					WithField("args", argWoFlags).
					Error("failed to pre run evaluation command")
				continue
			}
			err = cmd.RunE(cmd, argWoFlags)
			if err != nil {
				log.WithError(err).
					WithField("command", cmd.Name()).
					WithField("args", argWoFlags).
					Error("failed to run evaluation command")
				continue
			}
		}

		return nil
	},
}
