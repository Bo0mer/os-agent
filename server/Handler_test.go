package server_test

import (
	. "github.com/Bo0mer/os-agent/server"
	. "github.com/Bo0mer/os-agent/server/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {

	var handlerFunc HandlerFunc
	var handler Handler
	var argRequest Request
	var argResponse Response

	BeforeEach(func() {
		argRequest = nil
		argResponse = nil
		handlerFunc = func(req Request, resp Response) {
			argRequest = req
			argResponse = resp
		}

		handler = NewHandler("POST", "/this/is/sparta", handlerFunc)
	})

	It("has correct binding", func() {
		expectedBinding := Binding{
			Method: "POST",
			Path:   "/this/is/sparta",
		}

		Expect(handler.Binding()).To(Equal(expectedBinding))
	})

	It("calls handleFunc on Handle", func() {
		req := new(FakeRequest)
		resp := new(FakeResponse)
		handler.Handle(req, resp)
		Expect(argRequest).To(Equal(req))
		Expect(argResponse).To(Equal(resp))
	})

})
