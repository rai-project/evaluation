package writer

type Options struct {
	FilterKernelNames []string
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

func NewOptions(opts ...Option) Options {
	res := &Options{}

	for _, o := range opts {
		o(res)
	}

	return *res
}
