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
		var err error

		cmds := []cobra.Command{
			latencyCmd,
			layersCmd,
			layersTreeCmd,
			cudaLaunchCmd,
			eventflowCmd,
			durationCmd,
		}

		for _, cmd := range cmds {
			log.WithField("command", cmd.Name()).Debug("running command")
			argWoFlags := cmd.Flags().Args()
			err = cmd.PreRunE(cmd, argWoFlags)
			if err != nil {
				log.WithError(err).
					WithField("command", cmd.Name()).
					WithField("args", argWoFlags).
					Error("failed to pre run command")
				continue
			}
			err = cmd.RunE(cmd, argWoFlags)
			if err != nil {
				log.WithError(err).
					WithField("command", cmd.Name()).
					WithField("args", argWoFlags).
					Error("failed to run command")
				continue
			}
		}

		return nil
	},
}
