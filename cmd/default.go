package cmd

import "github.com/spf13/cobra"

var (
	defaultStepTraceDatabaseName = "carml_step_trace"
	defaultFullTraceDatabaseName = "carml_full_trace"
	defaultAccuracyDatabaseName  = "carml_accuracy"
	defaultDatabaseName          = map[string]string{
		"duration":    defaultStepTraceDatabaseName,
		"latency":     defaultStepTraceDatabaseName,
		"eventflow":   defaultFullTraceDatabaseName,
		"layers":      defaultFullTraceDatabaseName,
		"layer_tree":  defaultFullTraceDatabaseName,
		"cuda_launch": defaultFullTraceDatabaseName,
		"accuracy":    defaultAccuracyDatabaseName,
	}
	AllCmds = []*cobra.Command{
		latencyCmd,
		layersCmd,
		layersTreeCmd,
		//	cudaLaunchCmd,
		//	eventflowCmd,
		durationCmd,
		accuracyCmd,
	}
)
