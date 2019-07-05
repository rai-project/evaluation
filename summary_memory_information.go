package evaluation

import (
	"regexp"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
)

//easyjson:json
type GPUMemInformation struct {
	GPUID int `json:"gpuid,omitempty"`

	StartUsed  int64 `json:"start_used,omitempty"`
	StartFree  int64 `json:"start_free,omitempty"`
	StartTotal int64 `json:"start_total,omitempty"`

	FinishUsed  int64 `json:"finish_used,omitempty"`
	FinishFree  int64 `json:"finish_free,omitempty"`
	FinishTotal int64 `json:"finish_total,omitempty"`
}

//easyjson:json
type SystemMemoryInformation struct {
	StartAvailable int64 `json:"start_available,omitempty"`
	StartFree      int64 `json:"start_free,omitempty"`
	StartTotal     int64 `json:"start_total,omitempty"`

	FinishAvailable int64 `json:"finish_available,omitempty"`
	FinishFree      int64 `json:"finish_free,omitempty"`
	FinishTotal     int64 `json:"finish_total,omitempty"`
}

//easyjson:json
type RuntimeMemoryInformation struct {
	StartAlloc      int64 `json:"start_alloc,omitempty"`
	StartHeapAlloc  int64 `json:"start_heap_alloc,omitempty"`
	StartHeapSys    int64 `json:"start_heap_sys,omitempty"`
	StartTotalAlloc int64 `json:"start_total_alloc,omitempty"`

	FinishAlloc      int64 `json:"finish_alloc,omitempty"`
	FinishHeapAlloc  int64 `json:"finish_heap_alloc,omitempty"`
	FinishHeapSys    int64 `json:"finish_heap_sys,omitempty"`
	FinishTotalAlloc int64 `json:"finish_total_alloc,omitempty"`
}

//easyjson:json
type MemoryInformation struct {
	GPU     []GPUMemInformation      `json:"gpu,omitempty"`
	Runtime RuntimeMemoryInformation `json:"runtime,omitempty"`
	System  SystemMemoryInformation  `json:"system,omitempty"`
}

//easyjson:json
type SummaryMemoryInformation struct {
	SummaryBase        `json:",inline"`
	MemoryInformations []MemoryInformation `json:"memory_informations,omitempty"`
}

func (p Performance) MemoryInformationSummary(e Evaluation) (*SummaryMemoryInformation, error) {
	cPredictSpans := p.Spans().FilterByOperationName("c_predict", tracer)

	return &SummaryMemoryInformation{
		SummaryBase:        e.summaryBase(),
		MemoryInformations: cPredictSpans.MemoryInformation(),
	}, nil
}

func (spns Spans) MemoryInformation() []MemoryInformation {
	res := []MemoryInformation{}
	for _, s := range spns {
		info := memoryInformationFromSpan(s)
		if info == nil {
			continue
		}
		res = append(res, *info)
	}
	return res
}

var gpuIdSelectorRe = regexp.MustCompile(`^start_gpu\[(\d+)\]_.*`)

func memoryInformationFromSpan(span model.Span) *MemoryInformation {
	if len(span.Logs) == 0 {
		return nil
	}
	logs := span.Logs
	memInfo := MemoryInformation{}

	getGPUId := func(str string) int {
		for _, match := range gpuIdSelectorRe.FindAllString(str, -1) {
			idx := cast.ToInt(match)
			tbl := memInfo.GPU
			if len(memInfo.GPU) < idx {
				tbl = make([]GPUMemInformation, idx+1)
			}
			if len(memInfo.GPU) != 0 {
				copy(tbl, memInfo.GPU)
			}
			memInfo.GPU[idx].GPUID = idx
			memInfo.GPU = tbl
			return idx
		}
		pp.Println("not found...")
		return 0
	}

	for _, lg := range logs {
		for _, f := range lg.Fields {
			if strings.HasPrefix(f.Key, "start_mem_alloc") {
				memInfo.Runtime.StartAlloc = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "start_mem_heap_alloc") {
				memInfo.Runtime.StartHeapAlloc = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "start_mem_heap_sys") {
				memInfo.Runtime.StartHeapSys = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "start_mem_total_alloc") {
				memInfo.Runtime.StartTotalAlloc = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "start_mem_sys_available") {
				memInfo.System.StartAvailable = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "start_mem_sys_free") {
				memInfo.System.StartFree = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "start_mem_sys_total") {
				memInfo.System.StartTotal = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "start_gpu[") && strings.HasSuffix(f.Key, "]_mem_free") {
				gpuId := getGPUId(f.Key)
				memInfo.GPU[gpuId].StartFree = cast.ToInt64(f.Value)
			}
			if strings.HasPrefix(f.Key, "start_gpu[") && strings.HasSuffix(f.Key, "]_mem_total") {
				gpuId := getGPUId(f.Key)
				memInfo.GPU[gpuId].StartTotal = cast.ToInt64(f.Value)
			}
			if strings.HasPrefix(f.Key, "start_gpu[") && strings.HasSuffix(f.Key, "]_mem_used") {
				gpuId := getGPUId(f.Key)
				memInfo.GPU[gpuId].StartUsed = cast.ToInt64(f.Value)
			}

			if strings.HasPrefix(f.Key, "finish_mem_alloc") {
				memInfo.Runtime.FinishAlloc = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "finish_mem_heap_alloc") {
				memInfo.Runtime.FinishHeapAlloc = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "finish_mem_heap_sys") {
				memInfo.Runtime.FinishHeapSys = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "finish_mem_total_alloc") {
				memInfo.Runtime.FinishTotalAlloc = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "finish_mem_sys_available") {
				memInfo.System.FinishAvailable = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "finish_mem_sys_free") {
				memInfo.System.FinishFree = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "finish_mem_sys_total") {
				memInfo.System.FinishTotal = cast.ToInt64(f.Value)
				continue
			}
			if strings.HasPrefix(f.Key, "finish_gpu[") && strings.HasSuffix(f.Key, "]_mem_free") {
				gpuId := getGPUId(f.Key)
				memInfo.GPU[gpuId].FinishFree = cast.ToInt64(f.Value)
			}
			if strings.HasPrefix(f.Key, "finish_gpu[") && strings.HasSuffix(f.Key, "]_mem_total") {
				gpuId := getGPUId(f.Key)
				memInfo.GPU[gpuId].FinishTotal = cast.ToInt64(f.Value)
			}
			if strings.HasPrefix(f.Key, "finish_gpu[") && strings.HasSuffix(f.Key, "]_mem_used") {
				gpuId := getGPUId(f.Key)
				memInfo.GPU[gpuId].FinishUsed = cast.ToInt64(f.Value)
			}

		}
	}

	if isZero(memInfo) {
		return nil
	}

	return &memInfo
}
