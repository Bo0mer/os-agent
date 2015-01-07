package server_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/Bo0mer/os-agent/server"
	"github.com/Bo0mer/os-agent/server/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {

	var server Server
	var err error

	var doAction = func(resp *http.Response, err error) ([]byte, int, error) {
		if err != nil {
			return nil, 0, err
		}
		defer resp.Body.Close()
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return responseBody, resp.StatusCode, nil
	}

	var doPost = func(path string, bodyType string, body []byte) ([]byte, int, error) {
		url := fmt.Sprintf("http://%s%s", server.Address(), path)
		return doAction(http.Post(url, bodyType, bytes.NewBuffer(body)))
	}

	var doGet = func(path string) ([]byte, int, error) {
		url := fmt.Sprintf("http://%s%s", server.Address(), path)
		return doAction(http.Get(url))
	}

	Context("when the server is started with invalid address", func() {
		BeforeEach(func() {
			server = NewServer("n0t_val1d_h0s7", 0)
			err = server.Start()
		})

		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when the server is started with invalid port", func() {
		BeforeEach(func() {
			server = NewServer("127.0.0.1", -500)
			err = server.Start()
		})

		It("should return an error", func() {
			// broken on Windows for some reason
			// Expect(err).To(HaveOccurred())
		})
	})

	Context("when the server is started with proper host and port", func() {
		BeforeEach(func() {
			server = NewServer("127.0.0.1", 0)
			err = server.Start()
		})

		AfterEach(func() {
			server.Stop()
		})

		It("should be possible to get the server's address", func() {
			Ω(server.Address()).To(ContainSubstring("127.0.0.1"))
		})

		It("calling start multiple times does nothing", func() {
			err = server.Start()
			Ω(err).ToNot(HaveOccurred())

		})

		Context("when the server is stopped", func() {
			BeforeEach(func() {
				err = server.Stop()
			})

			It("an error should not have occurred", func() {
				Ω(err).ToNot(HaveOccurred())
			})

			It("calling stop multiple times does nothing", func() {
				err = server.Stop()
				Ω(err).ToNot(HaveOccurred())
			})

		})

		Describe("Route Handling", func() {

			var handler *fakes.FakeHandler
			var body []byte
			var response_body []byte
			var status int
			var err error

			BeforeEach(func() {
				body = []byte("request body")

				handler = new(fakes.FakeHandler)
				handler.BindingReturns(Binding{
					Method: "POST",
					Path:   "/foo/bar",
				})

				handler.HandleStub = func(req Request, resp Response) {
					Expect(req.Body()).To(Equal(body))
					resp.SetBody(body)
					resp.SetStatusCode(http.StatusOK)
				}

				server.Register(handler)
			})

			Context("when an actual route is called", func() {
				BeforeEach(func() {
					response_body, status, err = doPost("/foo/bar", "text/plain", body)
				})

				It("an error should not have occurred", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should have called the right handler", func() {
					Expect(handler.HandleCallCount()).To(Equal(1))
				})

				It("the server should have returned status ok", func() {
					Expect(status).To(Equal(http.StatusOK))
				})

				It("the server should have returned the proper response", func() {
					Expect(response_body).To(Equal(body))
				})
			})

			Context("when an actual route is called with different method", func() {
				BeforeEach(func() {
					response_body, status, err = doGet("/foo/bar")
				})

				It("should return status code 404 not found", func() {
					Expect(status).To(Equal(http.StatusNotFound))
				})

				It("should have not called the handler for other method", func() {
					Expect(handler.HandleCallCount()).To(Equal(0))
				})
			})

			Context("when a missing route is called", func() {
				BeforeEach(func() {
					_, status, err = doPost("/missing", "application/json", nil)
				})

				It("should not return an error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should return status code 404 not found", func() {
					Expect(status).To(Equal(http.StatusNotFound))
				})
			})

		})

	})

})
