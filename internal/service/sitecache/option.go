package sitecache

import "io"

// Option configures the sitecache instance..
type Option interface {
	HydratorOption
	RendererOption
}

type option struct {
	HydratorOption
	RendererOption
}

// WithOutput sets the output writer.
func WithOutput(w io.Writer) Option {
	return option{
		HydratorOption: hydratorOptionFunc(func(r *Hydrator) {
			r.output = w
		}),
		RendererOption: rendererOptionFunc(func(r *Renderder) {
			r.output = w
		}),
	}
}
