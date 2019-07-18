package writer

import (
	"strings"

	"github.com/getlantern/deepcopy"
)

type Options struct {
	FilterKernelNames []string
	ShowSummaryBase   bool
	Formats           []string
}

type Option func(*Options)

func FilterKernelNames(kernels []string) Option {
	return func(w *Options) {
		if kernels == nil {
			kernels = []string{}
		}
		w.FilterKernelNames = kernels
	}
}

func ShowSummaryBase(b bool) Option {
	return func(w *Options) {
		w.ShowSummaryBase = b
	}
}

func Format(f string) Option {
	return func(w *Options) {
		f := strings.ToLower(f)
		w.Formats = strings.Split(f, ",")
	}
}

func Formats(f string) Option {
	return func(w *Options) {
		f := strings.ToLower(f)
		w.Formats = strings.Split(f, ",")
	}
}

func FromOptions(os Options) Option {
	return func(w *Options) {
		err := deepcopy.Copy(w, os)
		if err != nil {
			panic(err)
		}
	}
}

func NewOptions(opts ...Option) Options {
	res := &Options{}

	for _, o := range opts {
		o(res)
	}

	return *res
}
