package api_test

import (
	"context"
	"github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/solo-io/ext-auth-plugins/api"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var _ = Describe("api has no errors", func() {

	It("can compile everything", func() {
		var pluginImpl api.ExtAuthPlugin = &pluginImpl{}
		_, err := pluginImpl.NewConfigInstance(context.Background())
		Expect(err).NotTo(HaveOccurred())

		svc, err := pluginImpl.GetAuthService(context.Background(), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())
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

func (serviceImpl) Authorize(ctx context.Context, request *v2.CheckRequest) (*api.AuthorizationResponse, error) {
	return nil, nil
}
