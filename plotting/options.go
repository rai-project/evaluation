package plotting

import (
	"path/filepath"

	"github.com/AlekSi/pointer"
	"github.com/mitchellh/go-homedir"
)

type Options struct {
	baseDir          string
	machineHostName  string
	frameworkName    string
	frameworkVersion string
	modelName        string
	modelVersion     string
	batchSize        int
	useGPU           *bool
	ignoreReadErrors bool
}

type OptionModifier func(o *Options)
type OptionModifiers struct{}

var Option = OptionModifiers{}

func (o OptionModifiers) MachineHostName(hostName string) OptionModifier {
	return func(o *Options) {
		o.machineHostName = hostName
	}
}

func (o OptionModifiers) BaseDir(dir string) OptionModifier {
	return func(o *Options) {
		o.baseDir = dir
	}
}

func (o OptionModifiers) FrameworkName(s string) OptionModifier {
	return func(o *Options) {
		o.frameworkName = s
	}
}

func (o OptionModifiers) FrameworkVersion(s string) OptionModifier {
	return func(o *Options) {
		o.frameworkVersion = s
	}
}

func (o OptionModifiers) ModelName(s string) OptionModifier {
	return func(o *Options) {
		o.modelName = s
	}
}

func (o OptionModifiers) ModelVersion(s string) OptionModifier {
	return func(o *Options) {
		o.modelVersion = s
	}
}

func (o OptionModifiers) BatchSize(val int) OptionModifier {
	return func(o *Options) {
		o.batchSize = val
	}
}

func (o OptionModifiers) UseGPU(val bool) OptionModifier {
	return func(o *Options) {
		o.useGPU = &val
	}
}

func (o OptionModifiers) IgnoreReadErrors(val bool) OptionModifier {
	return func(o *Options) {
		o.ignoreReadErrors = val
	}
}

func NewOptions(os ...OptionModifier) *Options {
	home, _ := homedir.Dir()
	opts := &Options{
		baseDir:          filepath.Join(home, "experiments"),
		machineHostName:  "ip-172-31-20-197",
		frameworkName:    "TensorFlow",
		frameworkVersion: "1.12",
		modelName:        "BVLC_AlexNet_Caffe",
		modelVersion:     "1.0",
		batchSize:        1,
		useGPU:           pointer.ToBool(true),
		ignoreReadErrors: false,
	}
	for _, o := range os {
		o(opts)
	}
	return opts
}
