package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	kernelNameFilterString string
	kernelNameFilterList   = []string{}
)

var gpuKernelInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{},
	Short:   "Get evaluation gpu kernel information from system library traces in a database. Specify model name as `all` to list information of all the models.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		err := rootSetup()
		if err != nil {
			return err
		}
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "cuda_kernel_launch")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}

		if kernelNameFilterString != "" {
			kernelNameFilterList = strings.Split(kernelNameFilterString, ",")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary0, err := evals.SummaryGPUKernelLayerInformations(performanceCollection)
			if err != nil {
				return err
			}

			if sortOutput || topLayers != -1 {
				sort.Sort(summary0)
				if topLayers != -1 {
					if topLayers >= len(summary0) {
						topLayers = len(summary0)
					}
					summary0 = summary0[:topLayers]
				}
				for ii := range summary0 {
					kernelInfo := summary0[ii]
					sort.Sort(kernelInfo)
					summary0[ii] = kernelInfo
				}
			}

			writer := NewWriter(summary0)
			defer writer.Close()
			for _, elem := range summary0 {
				writer.Rows(elem)
			}
			return nil
		}

		return forallmodels(run)
	},
}
