package server

type Handler interface {
	Binding() Binding
	Handle(Request, Response)
}

type HandlerFunc func(Request, Response)

type handler struct {
	binding     Binding
	handlerFunc HandlerFunc
}

func NewHandler(method, path string, handlerFunc HandlerFunc) Handler {
	binding := Binding{
		Method: method,
		Path:   path,
	}
	return &handler{
		binding:     binding,
		handlerFunc: handlerFunc,
	}
}

func (h *handler) Binding() Binding {
	return h.binding
}

func (h *handler) Handle(request Request, response Response) {
	h.handlerFunc(request, response)
}
