package server 

import "net/http"

type serveMux struct {
	*http.ServeMux
}

func NewServeMux() *serveMux {
	return &serveMux{http.NewServeMux()}
}


func (m *serveMux) Routes() {
	m.Handle("GET /", m.testHandler())
}
