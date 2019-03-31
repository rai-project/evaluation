package metrics

import (
	"sync"

	"github.com/rai-project/dlframework"
)

type FeatureCompareFunction func(featA, featB *dlframework.Feature) float64

type featureCompareRegistryMap struct {
	fs map[string]FeatureCompareFunction
	sync.Mutex
}

var featureCompareRegistry = featureCompareRegistryMap{
	fs: map[string]FeatureCompareFunction{},
}

func RegisterFeatureCompareFunction(name string, f FeatureCompareFunction) {
	featureCompareRegistry.Lock()
	featureCompareRegistry.fs[name] = f
	featureCompareRegistry.Unlock()
}
