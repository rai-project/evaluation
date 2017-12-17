package evaluation

import (
	"regexp"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
)

type GPUMemInformation struct {
	GPUID int

	StartUsed  int64
	StartFree  int64
	StartTotal int64

	FinishUsed  int64
	FinishFree  int64
	FinishTotal int64
}

type SystemMemoryInformation struct {
	StartAvailable int64
	StartFree      int64
	StartTotal     int64

	FinishAvailable int64
	FinishFree      int64
	FinishTotal     int64
}

type RuntimeMemoryInformation struct {
	StartAlloc      int64
	StartHeapAlloc  int64
	StartHeapSys    int64
	StartTotalAlloc int64

	FinishAlloc      int64
	FinishHeapAlloc  int64
	FinishHeapSys    int64
	FinishTotalAlloc int64
}

type MemoryInformation struct {
	GPU     []GPUMemInformation
	Runtime RuntimeMemoryInformation
	System  SystemMemoryInformation
}

type SummaryMemoryInformation struct {
	SummaryBase
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	HostName            string
	MemoryInformations  []MemoryInformation
}

func (p Performance) MemoryInformationSummary(e Evaluation) (*SummaryMemoryInformation, error) {
	spans := p.Spans().FilterByOperationName("predict")

	return &SummaryMemoryInformation{
		SummaryBase:         e.summaryBase(),
		MachineArchitecture: e.MachineArchitecture,
		UsingGPU:            e.UsingGPU,
		BatchSize:           e.BatchSize,
		HostName:            e.Hostname,
		MemoryInformations:  spans.MemoryInformation(),
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
