package cmd

import "github.com/spf13/cobra"

var (
	defaultModelTraceDatabaseName         = "carml_model_trace"
	defaultFrameworkTraceDatabaseName     = "carml_framework_trace"
	defaultSystemLibraryTraceDatabaseName = "carml_system_library_trace"
	defaultFullTraceDatabaseName          = "carml_full_trace"
	defaultAccuracyDatabaseName           = "carml_accuracy"
	defaultDatabaseName                   = map[string]string{
		"model":       defaultModelTraceDatabaseName,
		"layer":       defaultFrameworkTraceDatabaseName,
		"cuda_kernel": defaultSystemLibraryTraceDatabaseName,
		"accuracy":    defaultAccuracyDatabaseName,
	}
	AllCmds = []*cobra.Command{
		modelCmd,
		layerCmd,
		cudaKernelCmd,
		eventflowCmd,
		accuracyCmd,
	}
)
