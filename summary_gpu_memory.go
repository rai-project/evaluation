package evaluation

import (
	"strings"
)

func isCUDAMemOp(name string) bool {
	name = strings.ToLower(name)
	switch name {
	case "cudamalloc",
		"cudafree",
		"cuda_memcpy",
		"cudamemcpy",
		"cudamallochost",
		"cudafreehost",
		"cuda_memcpy_dev":
		return true
	}
	return false
}
