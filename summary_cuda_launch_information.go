package evaluation

type KernelLaunchInformation struct {
	Name       string    `json:"name,omitempty"`
	Parameters []string  `json:"parameters,omitempty"`
	Durations  []float64 `json:"durations,omitempty"`
}

type SummaryCUDALaunchInformation struct {
	SummaryBase              `json:",inline"`
	MachineArchitecture      string                    `json:"machine_architecture,omitempty"`
	UsingGPU                 bool                      `json:"using_gpu,omitempty"`
	BatchSize                int                       `json:"batch_size,omitempty"`
	HostName                 string                    `json:"host_name,omitempty"`
	KernelLaunchInformations []KernelLaunchInformation `json:"kernel_launch_information,omitempty"`
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
