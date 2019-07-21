package evaluation

import (
	json "encoding/json"
	"strings"

	model "github.com/uber/jaeger/model/json"
)

type TensorFlowTagOutput struct {
	TensorDescription TensorDescription `json:"tensor_description"`
}
type Dim struct {
	Size int `json:"size"`
}
type Shape struct {
	Dim []Dim `json:"dim"`
}
type AllocationDescription struct {
	RequestedBytes int    `json:"requested_bytes"`
	AllocatedBytes int    `json:"allocated_bytes"`
	AllocatorName  string `json:"allocator_name"`
	AllocationID   int    `json:"allocation_id"`
	Ptr            int64  `json:"ptr"`
}
type TensorDescription struct {
	Dtype                 int                   `json:"dtype"`
	Shape                 Shape                 `json:"shape"`
	AllocationDescription AllocationDescription `json:"allocation_description"`
}

func getAllocationDescription(span model.Span) AllocationDescription {
	ret := AllocationDescription{}
	output, err := getTagValueAsString(span, "output")
	if err != nil {
		log.WithError(err).Info("fail to get output value in the tags")
		return ret
	}
	if output == "" {
		return ret
	}
	output = strings.Replace(output, "\\", "", -1)

	var result []TensorFlowTagOutput
	json.Unmarshal([]byte(output), &result)

	if len(result) == 0 {
		return ret
	}
	ret = result[0].TensorDescription.AllocationDescription
	return ret
}

type TensorFlowAllocatorMemoryUsed struct {
	AllocatorName       string              `json:"allocator_name"`
	TotalBytes          int                 `json:"total_bytes"`
	PeakBytes           int                 `json:"peak_bytes"`
	LiveBytes           int                 `json:"live_bytes"`
	AllocationRecords   []AllocationRecords `json:"allocation_records"`
	AllocatorBytesInUse int                 `json:"allocator_bytes_in_use"`
}
type AllocationRecords struct {
	AllocMicros int64 `json:"alloc_micros"`
	AllocBytes  int   `json:"alloc_bytes"`
}

func getTensorFlowAllocatorMemoryUsed(span model.Span) (TensorFlowAllocatorMemoryUsed, bool) {
	ret := TensorFlowAllocatorMemoryUsed{}
	output, err := getTagValueAsString(span, "memory")
	if err != nil {
		return ret, false
	}
	if output == "" {
		return ret, false
	}
	output = strings.Replace(output, "\\", "", -1)

	var result []TensorFlowAllocatorMemoryUsed
	json.Unmarshal([]byte(output), &result)

	if len(result) == 0 {
		return ret, false
	}
	ret = result[0]
	return ret, true
}
