package cmd

import (
	"github.com/spf13/cobra"
)

var (
	topKernels int
)

var gpuKernelCmd = &cobra.Command{
	Use: "gpu_kernel",
	Aliases: []string{
		"cuda_kernel",
		"kernel",
		"gpu_kernel",
	},
	Short: "Get evaluation model layer analysis from framework traces in a database",
}

func init() {
	gpuKernelCmd.PersistentFlags().StringVar(&kernelNameFilterString, "kernel_names", "", "filter out certain kernel (input must be mangled and is comma seperated)")
	gpuKernelCmd.PersistentFlags().IntVar(&topKernels, "top_kernels", -1, "consider only the top k kernel ranked by duration")

	gpuKernelCmd.AddCommand(gpuKernelInfoCmd)
	gpuKernelCmd.AddCommand(gpuKernelNameAggreCmd)
	gpuKernelCmd.AddCommand(gpuKernelModelAggreCmd)
	gpuKernelCmd.AddCommand(gpuKernelLayerAggreCmd)
	gpuKernelCmd.AddCommand(gpuKernelLayerFlopsCmd)
	gpuKernelCmd.AddCommand(gpuKernelLayerDramReadCmd)
	gpuKernelCmd.AddCommand(gpuKernelLayerDramWriteCmd)
	gpuKernelCmd.AddCommand(gpuKernelLayerDurationCmd)
	gpuKernelCmd.AddCommand(gpuKernelLayerGPUCPUCmd)
}
