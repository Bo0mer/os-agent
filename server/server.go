package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type Server interface {
	Start() error
	Stop() error
	Address() string
	Register(handler Handler)
}

type server struct {
	host         string
	port         int
	listener     net.Listener
	requestGroup sync.WaitGroup
	handlers     []Handler
}

type Request interface {
	Body() []byte
	ParamValues(string) ([]string, bool)
}

type request struct {
	body   []byte
	params map[string][]string
}

func (r *request) Body() []byte {
	return r.body
}

func (r *request) ParamValues(name string) (values []string, found bool) {
	values, found = r.params[name]
	return
}

type Response interface {
	SetBody([]byte)
	SetStatusCode(int)
	StatusCode() int
}

type response struct {
	body       []byte
	statusCode int
}

func (r *response) SetBody(data []byte) {
	r.body = data
}

func (r *response) SetStatusCode(statusCode int) {
	r.statusCode = statusCode
}

func (r *response) StatusCode() int {
	return r.statusCode
}

func NewServer(host string, port int) Server {
	return &server{
		host:     host,
		port:     port,
		handlers: []Handler{},
	}
}

func (s *server) Start() error {
	if s.listener != nil {
		return nil
	}
	var err error
	url := fmt.Sprintf("%s:%d", s.host, s.port)
	s.listener, err = net.Listen("tcp", url)
	if err != nil {
		return err
	}
	go http.Serve(s.listener, http.HandlerFunc(s.serveHTTP))
	return nil
}

func (s *server) Stop() error {
	if s.listener == nil {
		return nil
	}

	// wait for active requests to finish
	s.requestGroup.Wait()
	s.listener.Close()
	s.listener = nil
	return nil
}

func (s *server) Address() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

func (s *server) Register(handler Handler) {
	s.handlers = append(s.handlers, handler)
}

func (s *server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.requestGroup.Add(1)
	defer s.requestGroup.Done()

	// find handler to handle the reuqest
	handler, found := s.findHandler(r)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	request, err := s.createRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// otherwise generate response and pass it to handler
	response := &response{
		statusCode: http.StatusOK,
	}

	handler.Handle(request, response)

	w.WriteHeader(response.StatusCode())
	bytes.NewBuffer(response.body).WriteTo(w)
}

func (s *server) createRequest(r *http.Request) (*request, error) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	req := &request{
		body:   content,
		params: map[string][]string(r.URL.Query()),
	}
	return req, err
}

func (s *server) findHandler(r *http.Request) (Handler, bool) {
	binding := Binding{
		Method: r.Method,
		Path:   r.URL.Path,
	}

	for _, handler := range s.handlers {
		if handler.Binding() == binding {
			return handler, true
		}
	}
	return nil, false
}
