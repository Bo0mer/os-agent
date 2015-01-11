package masterclient_test

import (
	"net/http"

	. "github.com/Bo0mer/os-agent/masterclient"
	"github.com/Bo0mer/os-agent/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("MasterClient", func() {

	var fakeServer *ghttp.Server
	var self model.Slave
	var masterClient MasterClient

	BeforeEach(func() {
		fakeServer = ghttp.NewServer()
		self = model.Slave{
			Id:   "self-id",
			Host: "127.0.0.1",
			Port: 80,
		}
		masterClient = NewMasterClient(fakeServer.URL(), self)
	})

	Describe("Register", func() {

		Context("when the server returns non-2XX status code", func() {
			BeforeEach(func() {
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/register"),
					ghttp.RespondWithJSONEncoded(http.StatusNotFound, []byte{})),
				)
			})

			It("returns an error", func() {
				err := masterClient.Register()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the server returns OK", func() {
			BeforeEach(func() {
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/register"),
					ghttp.VerifyJSONRepresenting(self),
					ghttp.RespondWithJSONEncoded(http.StatusOK, []byte{})),
				)
			})

			It("returns no error", func() {
				err := masterClient.Register()
				Expect(err).ToNot(HaveOccurred())
			})

		})

	})

})
