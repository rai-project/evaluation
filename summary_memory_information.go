package evaluation

type GPUMemInformation struct {
}

type SystemMemoryInformation struct {
}

type MemoryInformation struct {
	GPU    GPUMemInformation
	System SystemMemoryInformation
}

type SummaryCUDALaunchInformation struct {
	SummaryBase
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	HostName            string
	MemoryInformations  []MemoryInformation
}
