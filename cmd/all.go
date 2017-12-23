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
		for _, pcmd := range allCmds {
			if !rootFlags.Changed("database_name") {
				switch cmd.Name() {
				case latencyCmd.Name(),
					durationCmd.Name():
					databaseName = "carml_step_trace"
				case layersCmd.Name(),
					layersTreeCmd.Name(),
					cudaLaunchCmd.Name(),
					eventflowCmd.Name():
					databaseName = "carml_full_trace"
				}
			}
			pcmd.Flags().AddFlagSet(rootFlags)
			log.WithField("command", pcmd.Name()).
				WithField("args", args).
				Info("running evaluation command")
			pp.Println("running evaluation command ", args, "  ", pcmd.Name())
			pcmd.SilenceUsage = true
			rootSetup()
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
