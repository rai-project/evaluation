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

func (p Performance) CUDALaunchInformationSummary(e Evaluation) (*SummaryPredictDurationInformation, error) {
	spans := p.Spans().FilterByOperationName("launch_kernel")

	return &SummaryPredictDurationInformation{
		SummaryBase:         e.summaryBase(),
		MachineArchitecture: e.MachineArchitecture,
		UsingGPU:            e.UsingGPU,
		BatchSize:           e.BatchSize,
		HostName:            e.Hostname,
		Durations:           spans.Duration(),
	}, nil
}
