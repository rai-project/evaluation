package evaluation

type SummaryPredictDurationInformation struct {
	SummaryBase
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	HostName            string
	Duration            float64 // in nano seconds
	Latency             float64 // in nano seconds
	Throughput          float64
}
