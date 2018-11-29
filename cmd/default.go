package cmd

import "github.com/spf13/cobra"

var (
	defaultModelTraceDatabaseName = "carml_model_trace"
	defaultFullTraceDatabaseName  = "carml_full_trace"
	defaultAccuracyDatabaseName   = "carml_accuracy"
	defaultDatabaseName           = map[string]string{
		"duration":    defaultModelTraceDatabaseName,
		"latency":     defaultModelTraceDatabaseName,
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
		cudaLaunchCmd,
		eventflowCmd,
		durationCmd,
		accuracyCmd,
	}
)
