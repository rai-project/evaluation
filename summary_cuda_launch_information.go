package evaluation

type KernelLaunchInformation struct {
	Name       string
	Parameters []string
	Durations  []float64
}

type SummaryCUDALaunchInformation struct {
	SummaryBase
	MachineArchitecture      string
	UsingGPU                 bool
	BatchSize                int
	HostName                 string
	KernelLaunchInformations []KernelLaunchInformation
}
