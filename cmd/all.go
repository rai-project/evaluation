package cmd

import (
	"os"

	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
)

func runAll(cmd *cobra.Command, args []string) error {
	rootFlags := cmd.Parent().Flags()
	rootFlags.Parse(os.Args[2:])
	originalDatabaseName := databaseName
	for _, pcmd := range allCmds {
		//pcmd.SilenceUsage = true
		if originalDatabaseName == "" {
			switch pcmd.Name() {
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
		pargs := append([]string{pcmd.Name()}, os.Args[2:]...)
		pargs = append(pargs, "--database_name="+databaseName)
		cmd.Parent().SetArgs(pargs)
		pcmd.SetArgs(pargs)

		pcmd.Flags().AddFlagSet(rootFlags)
		pcmd.Flags().Parse(pargs)
		rootSetup()

		log.WithField("command", pcmd.Name()).
			WithField("args", pargs).
			Info("running evaluation command")
		//pp.Println("running evaluation command ", pargs, "  ", pcmd.Name())
		err := pcmd.PreRunE(pcmd, pargs)
		if err != nil {
			log.WithError(err).
				WithField("command", pcmd.Name()).
				WithField("args", pargs).
				Error("failed to pre run evaluation command")
			continue
		}
		err = pcmd.RunE(pcmd, pargs)
		if err != nil {
			log.WithError(err).
				WithField("command", pcmd.Name()).
				WithField("args", pargs).
				Error("failed to run evaluation command")
			continue
		}
	}

	return nil
}

var allCmd = &cobra.Command{
	Use: "all",
	Aliases: []string{
		"eval_all",
	},
	Short: "Get all evaluation information from CarML",
	RunE: func(cmd *cobra.Command, args []string) error {
		pp.Println("use the run_all.sh script")

		return nil
	},
}
