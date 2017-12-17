package evaluation

type LayerInformation struct {
	Name      string
	Durations []float64
}

type SummaryLayerInformation struct {
	SummaryBase
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	HostName            string
	LayerInformations   []LayerInformation
}
