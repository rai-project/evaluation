package evaluation

type KernelLaunchInformation struct {
	Name       string    `json:"name,omitempty"`
	Parameters []string  `json:"parameters,omitempty"`
	Durations  []float64 `json:"durations,omitempty"`
}

type SummaryCUDALaunchInformation struct {
	SummaryBase              `json:",inline"`
	KernelLaunchInformations []KernelLaunchInformation `json:"kernel_launch_information,omitempty"`
}

func (p Performance) CUDALaunchInformationSummary(e Evaluation) (*SummaryPredictDurationInformation, error) {
	spans := p.Spans().FilterByOperationName("launch_kernel")

	return &SummaryPredictDurationInformation{
		SummaryBase: e.summaryBase(),
		Durations:   spans.Duration(),
	}, nil
}
