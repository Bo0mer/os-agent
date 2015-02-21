package config_test

import (
	. "github.com/Bo0mer/os-agent/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var err error
	var config OSAgentConfig

	Context("when the config is there", func() {

		BeforeEach(func() {
			config, err = LoadConfig("config.yml")
		})

		It("should not have returned an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should have the OS Agent Id config", func() {
			Expect(config.Id).To(Equal("unique-id"))
		})

		It("should have the OS Agent host config", func() {
			Expect(config.Host).To(Equal("127.0.0.1"))
		})

		It("should have the OS Agent port config", func() {
			Expect(config.Port).To(Equal(8080))
		})

		It("should have the server configuration", func() {
			Expect(config.Server.Host).To(Equal("127.0.0.1"))
			Expect(config.Server.Port).To(Equal(8080))
		})

		It("should have the Auth configuration", func() {
			Expect(config.Server.Auth.Enabled).To(BeTrue())
			Expect(config.Server.Auth.User).To(Equal("admin"))
			Expect(config.Server.Auth.Password).To(Equal("s3cre7"))
		})

		It("should have the master configuration", func() {
			Expect(config.Master.URL).To(Equal("http://127.0.0.1:8081"))
		})

	})
})
