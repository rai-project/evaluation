package cmd

import (
	"github.com/spf13/cobra"
)

var (
	topKernels int
)

var cudaKernelCmd = &cobra.Command{
	Use: "cuda_kernel",
	Aliases: []string{
		"cuda",
		"gpu",
		"kernel",
		"kernels",
		"gpu_kernel",
		"gpu_kernels",
	},
	Short: "Get evaluation model layer analysis from framework traces in a database",
}

func init() {
	cudaKernelCmd.PersistentFlags().StringVar(&kernelNameFilterString, "kernel_names", "", "filter out certain kernel (input must be mangled and is comma seperated)")
	cudaKernelCmd.PersistentFlags().IntVar(&topKernels, "top_kernels", -1, "consider only the top k kernel ranked by duration")

	cudaKernelCmd.AddCommand(cudaKernelInfoCmd)
	cudaKernelCmd.AddCommand(cudaKernelDurationCmd)
}
