package util

import (
	"net/http"
)


type Session struct {
	w      http.ResponseWriter
	r      *http.Request
	params func(key string) string
	flags  map[string]string
}

func (s *Session) Init(w http.ResponseWriter, r *http.Request) {
	s.w = w
	s.r = r
	s.params = r.FormValue
}

func (s *Session) Params(args ...string) (result string) {
	l := len(args)
	var name, def string

	switch {
	case l == 0:
		return
	case l == 1:
		name = args[0]
	default:
		name = args[0]
		def = args[1]
	}

	if result = s.params(name); result == "" {
		result = def
	}

	return
}

func (s *Session) SetResponseHeader(m map[string]string) {
	for k, v := range m {
		s.w.Header().Set(k, v)
	}
}

func (s *Session) ResponseEnd(json []byte, err error) {
	if err == nil {
		s.w.Write(json)
	} else {
		s.w.Write([]byte(err.Error()))
	}
}
