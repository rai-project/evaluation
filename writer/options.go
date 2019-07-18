package writer

type Options struct {
	FilterKernelNames []string
	ShowSummaryBase   bool
	Format            string
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
		w.Format = f
	}
}

func NewOptions(opts ...Option) Options {
	res := &Options{}

	for _, o := range opts {
		o(res)
	}

	return *res
}
