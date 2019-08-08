package main

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-plugins/api"
	api2 "github.com/solo-io/ext-auth-plugins/example/header/api"
	"plugin"
)

var _ = Describe("Plugin", func() {

	It("can be loaded", func() {

		goPlugin, err := plugin.Open("RequiredHeader.so")
		Expect(err).NotTo(HaveOccurred())

		pluginStructPtr, err := goPlugin.Lookup("Plugin")
		Expect(err).NotTo(HaveOccurred())

		extAuthPlugin, ok := pluginStructPtr.(api.ExtAuthPlugin)
		Expect(ok).To(BeTrue())

		instance := extAuthPlugin.NewConfigInstance(context.TODO())

		typedInstance, ok := instance.(api2.RequiredHeaderPlugin)
		Expect(ok).To(BeTrue())

		Expect(typedInstance.RequiredHeader).To(BeEmpty())
	})
})
