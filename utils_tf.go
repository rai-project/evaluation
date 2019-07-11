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

func getAllocationBytes(span model.Span) int64 {
	output, err := getTagValueAsString(span, "output")
	if err != nil {
		log.WithError(err).Info("fail to get output value in the tags")
		return int64(0)
	}
	if output == "" {
		return int64(0)
	}
	output = strings.Replace(output, "\\", "", -1)

	var result []TensorFlowTagOutput
	json.Unmarshal([]byte(output), &result)

	if len(result) == 0 {
		return 0
	}

	return int64(result[0].TensorDescription.AllocationDescription.AllocatedBytes)
}
