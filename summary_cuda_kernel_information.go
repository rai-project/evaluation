package evaluation

import (
	"strings"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
)

type Metadata map[string]interface{}

type SummaryCUDAKernelInformation struct {
	Name               string     `json:"name,omitempty"`
	MangledName        string     `json:"mangled_name,omitempty"`
	Durations          []int64    `json:"durations,omitempty"`
	Tags               []Metadata `json:"tags,omitempty"`
	Logs               []Metadata `json:"logs,omitempty"`
	CorrelationId      int64      `json:"correlation_id,omitempty"`
	MeanDuration       float64    `json:"mean_duration,omitempty"`
	MeanFlops          float64    `json:"mean_flops,omitempty"`
	MeanDramReadBytes  float64    `json:"mean_dram_read_bytes,omitempty"`
	MeanDramWriteBytes float64    `json:"mean_dram_write_bytes,omitempty"`
}

type SummaryCUDAKernelInformations []SummaryCUDAKernelInformation

func (p SummaryCUDAKernelInformations) Len() int { return len(p) }
func (p SummaryCUDAKernelInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	xDuration := TrimmedMeanInt64Slice(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMeanInt64Slice(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p SummaryCUDAKernelInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func getMetaDataKeys(metadatas []Metadata) []string {
	if metadatas == nil || len(metadatas) == 0 {
		return []string{}
	}
	keyVisited := map[string]bool{}
	keys := []string{}
	for _, metadata := range metadatas {
		for key, _ := range metadata {
			if _, ok := keyVisited[key]; ok {
				continue
			}
			keys = append(keys, key)
			keyVisited[key] = true
		}
	}
	return keys
}

func (infos SummaryCUDAKernelInformations) GetKernelLogKeys() []string {
	kernelLogs := []Metadata{}
	for _, info := range infos {
		if len(info.Logs) == 0 {
			continue
		}
		kernelLogs = append(kernelLogs, info.Logs...)
	}
	return getMetaDataKeys(kernelLogs)
}

func (info SummaryCUDAKernelInformation) Header(opts ...writer.Option) []string {
	extraHeader := []string{
		"kernel_name",
		"kernel_mean_duration (us)",
		"kernel_mean_flops",
		"kernel_mean_dram_read_bytes",
		"kernel_mean_dram_write_bytes",
		"kernel_durations (us)",
	}
	kernelLogKeys := SummaryCUDAKernelInformations{info}.GetKernelLogKeys()
	if len(kernelLogKeys) != 0 {
		extraHeader = append(extraHeader, kernelLogKeys...)
	}
	return extraHeader
}

func (info SummaryCUDAKernelInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		info.Name,
		cast.ToString(info.MeanDuration),
		cast.ToString(info.MeanFlops),
		cast.ToString(info.MeanDramReadBytes),
		cast.ToString(info.MeanDramWriteBytes),
		strings.Join(int64SliceToStringSlice(info.Durations), DefaultDimiter),
	}
	kernelLogKeys := SummaryCUDAKernelInformations{info}.GetKernelLogKeys()
	for _, kernelLogKey := range kernelLogKeys {
		kernelLogs := []string{}
		for _, kernelLog := range info.Logs {
			for kernelLogKeyName, keryeLogValue := range kernelLog {
				if kernelLogKeyName == kernelLogKey {
					kernelLogs = append(kernelLogs, cast.ToString(keryeLogValue))
				}
			}
		}
		extra = append(extra, strings.Join(kernelLogs, DefaultDimiter))
	}
	return extra
}

func (k *SummaryCUDAKernelInformation) addLogs(spanLogs []model.Log) {
	if k.Logs == nil {
		k.Logs = []Metadata{}
	}
	logs := Metadata{}
	for _, v := range spanLogs {
		for _, f := range v.Fields {
			logs[f.Key] = f.Value
		}
	}
	if len(logs) == 0 {
		return
	}
	k.Logs = append(k.Logs, logs)
}

func (k *SummaryCUDAKernelInformation) addTags(spanTags []model.KeyValue) {
	if k.Tags == nil {
		k.Tags = []Metadata{}
	}
	tags := Metadata{}
	for _, v := range spanTags {
		tags[v.Key] = v.Value
	}
	if len(tags) == 0 {
		return
	}
	k.Tags = append(k.Tags, tags)
}

func GetMeanLogValue(info SummaryCUDAKernelInformation, name string, trimmedMeanFraction float64) float64 {
	if info.Logs == nil {
		info.Logs = []Metadata{}
	}
	kernelLogs := []int64{}
	for _, kernelLog := range info.Logs {
		for kernelLogKeyName, keryeLogValue := range kernelLog {
			if kernelLogKeyName == name {
				kernelLogs = append(kernelLogs, cast.ToInt64(keryeLogValue))
			}
		}
	}
	return TrimmedMeanInt64Slice(kernelLogs, trimmedMeanFraction)
}

func GPUKernelSpantoCUDAKernelInformation(span model.Span) SummaryCUDAKernelInformation {
	info := &SummaryCUDAKernelInformation{
		Name:          mustGetTagValueAsString(span, "kernel_name"),
		MangledName:   mustGetTagValueAsString(span, "name"),
		Tags:          []Metadata{},
		Logs:          []Metadata{},
		CorrelationId: mustGetTagValueAsInt64(span, "correlation_id"),
		Durations: []int64{
			cast.ToInt64(span.Duration),
		},
	}
	info.addTags(span.Tags)
	info.addLogs(span.Logs)
	return *info
}

func CUDALaunchSpantoCUDAKernelInformation(span model.Span) SummaryCUDAKernelInformation {
	kernelName := mustGetTagValueAsString(span, "kernel")
	info := &SummaryCUDAKernelInformation{
		Name:          demangleName(kernelName),
		MangledName:   kernelName,
		Tags:          []Metadata{},
		Logs:          []Metadata{},
		CorrelationId: mustGetTagValueAsInt64(span, "correlation_id"),
	}
	info.addTags(span.Tags)
	info.addLogs(span.Logs)
	return *info
}

func (es Evaluations) CUDAKernelInformationSummary(perfCol *PerformanceCollection) (SummaryCUDAKernelInformations, error) {
	summary := SummaryCUDAKernelInformations{}

	layerCUDAKernelInfos, err := es.LayerCUDAKernelInformationSummary(perfCol)
	if err != nil {
		return summary, err
	}
	for _, info := range layerCUDAKernelInfos {
		summary = append(summary, info.SummaryCUDAKernelInformations...)
	}
	return summary, nil
}
