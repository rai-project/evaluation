package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func runAll(cmd *cobra.Command, args []string) error {
	rootFlags := cmd.Parent().Flags()
	rootFlags.Parse(os.Args[2:])
	originalDatabaseName := databaseName
	for _, pcmd := range AllCmds {
		//pcmd.SilenceUsage = true
		if originalDatabaseName == "" {
			switch pcmd.Name() {
			case latencyCmd.Name(),
				durationCmd.Name():
				databaseName = "carml_model_trace"
			case layerCmd.Name(),
				cudaKernelCmd.Name(),
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
		err := rootSetup()
		if err != nil {
			return err
		}
		log.WithField("command", pcmd.Name()).
			WithField("args", pargs).
			Info("running evaluation command")
		//pp.Println("running evaluation command ", pargs, "  ", pcmd.Name())
		err = pcmd.PreRunE(pcmd, pargs)
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
	Short: "Get all evaluation information from MLModelScope",
	RunE: func(*cobra.Command, []string) error {

		buildFile, err := getBuildFile()
		if err != nil {
			return err
		}

		cmd := exec.Command("go", "build", buildFile)
		if err := cmd.Run(); err != nil {
			return err
		}

		exe := filepath.Join(sourcePath, "main")

		args := os.Args[2:]
		for _, pcmd := range AllCmds {
			pargs := append([]string{pcmd.Name()}, args...)
			log.WithField("command", pcmd.Name()).
				WithField("args", pargs).
				Info("running evaluation command")
			cmd = exec.Command(exe, pargs...)
			if err := cmd.Run(); err != nil {
				log.WithError(err).
					WithField("command", pcmd.Name()).
					WithField("args", pargs).
					Error("failed to pre run evaluation command")
				continue
			}
		}

		return nil
	},
}
