package cmd

import "github.com/spf13/cobra"

var (
	defaultModelTraceDatabaseName         = "carml_model_trace"
	defaultFrameworkTraceDatabaseName     = "carml_framework_trace"
	defaultSystemLibraryTraceDatabaseName = "carml_system_library_trace"
	defaultFullTraceDatabaseName          = "carml_full_trace"
	defaultAccuracyDatabaseName           = "carml_accuracy"
	defaultDatabaseName                   = map[string]string{
		"duration":    defaultModelTraceDatabaseName,
		"latency":     defaultModelTraceDatabaseName,
		"eventflow":   defaultFrameworkTraceDatabaseName,
		"layer":       defaultFrameworkTraceDatabaseName,
		"cuda_kernel": defaultSystemLibraryTraceDatabaseName,
		"accuracy":    defaultAccuracyDatabaseName,
	}
	AllCmds = []*cobra.Command{
		latencyCmd,
		layersCmd,
		cudaKernelCmd,
		eventflowCmd,
		durationCmd,
		accuracyCmd,
	}
)
