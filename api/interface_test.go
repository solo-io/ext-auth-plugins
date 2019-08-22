package api_test

import (
	"context"

	"github.com/solo-io/ext-auth-plugins/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("api has no errors", func() {

	It("can compile everything", func() {
		var pluginImpl api.ExtAuthPlugin = &pluginImpl{}
		_, err := pluginImpl.NewConfigInstance(context.Background())
		Expect(err).NotTo(HaveOccurred())

		svc, err := pluginImpl.GetAuthService(context.Background(), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())
	})

	It("can set state", func() {
		var a api.AuthorizationRequest
		a.SetState("test", "testState")
		Expect(a.GetState("test")).To(BeEquivalentTo("testState"))
		a.SetState("test", 123)
		Expect(a.GetState("test")).To(BeEquivalentTo(123))
	})
	It("should not crash when only get state is called", func() {
		var a api.AuthorizationRequest
		Expect(a.GetState("test")).To(BeNil())
	})
})

type pluginImpl struct{}

func (pluginImpl) NewConfigInstance(ctx context.Context) (configInstance interface{}, err error) {
	return nil, nil
}

func (pluginImpl) GetAuthService(ctx context.Context, configInstance interface{}) (api.AuthService, error) {
	return &serviceImpl{}, nil
}

type serviceImpl struct{}

func (serviceImpl) Start(ctx context.Context) error {
	return nil
}

func (serviceImpl) Authorize(ctx context.Context, request *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {
	return nil, nil
}
