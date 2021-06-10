package server

type RouterHandlerOption func(opts *RouterHandlerOptions)

type RouterHandlerOptions struct {
	ParametersIn                string
	Summary                     string
	Tags                        []string
	ResponseWrapper             SwaggerResponseWrapper
	ResponseWrapperDataNodeName string
}

type SwaggerResponseWrapper interface {
	SetData(v interface{})
}

func NewRouterHandlerOptions(options ...RouterHandlerOption) *RouterHandlerOptions {
	o := &RouterHandlerOptions{
		ParametersIn:                SwaggerParametersInBody,
		Tags:                        make([]string, 0),
		ResponseWrapperDataNodeName: "data",
	}
	for _, opt := range options {
		opt(o)
	}
	return o
}

func (s *RouterHandlerOptions) SetParametersIn(v string) *RouterHandlerOptions {
	s.ParametersIn = v
	return s
}

func NewSwaggerParametersInOption(v string) RouterHandlerOption {
	return func(opts *RouterHandlerOptions) {
		opts.ParametersIn = v
	}
}

func NewSwaggerSummaryOption(v string) RouterHandlerOption {
	return func(opts *RouterHandlerOptions) {
		opts.Summary = v
	}
}

func NewSwaggerTagsOption(v []string) RouterHandlerOption {
	return func(opts *RouterHandlerOptions) {
		opts.Tags = v
	}
}

func NewSwaggerResponseWrapper(v SwaggerResponseWrapper, dataNodeName string) RouterHandlerOption {
	return func(opts *RouterHandlerOptions) {
		opts.ResponseWrapper = v
		opts.ResponseWrapperDataNodeName = dataNodeName
	}
}
