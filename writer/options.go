package writer

type Options struct {
	filterKernelNames []string
}

type Option func(*Options)

func FilterKernelNames(kernels []string) Option {
	return func(w *Options) {
		w.filterKernelNames = kernels
	}
}

func NewOptions(opts ...Option) Options {
	res := &Options{}

	for _, o := range opts {
		o(res)
	}

	return *res
}
