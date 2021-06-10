package server

import "context"

type Context struct {
	context.Context
}

var ctxRequest struct{}

func (s Context) Request() *http.Request {
	v := s.Value(ctxRequest)
	if v == nil {
		return nil
	}
	return v.(*http.Request)
}
