package api

import (
	"context"
	"errors"
	"fmt"
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
)

var _ api.ExtAuthPlugin = new(RequiredHeaderPlugin)

type RequiredHeaderPlugin struct{}

type RequiredHeaderPluginConfig struct {
	RequiredHeader string
}

func (p *RequiredHeaderPlugin) NewConfigInstance(ctx context.Context) (interface{}, error) {
	return &RequiredHeaderPluginConfig{}, nil
}

func (p *RequiredHeaderPlugin) GetAuthService(ctx context.Context, configInstance interface{}) (api.AuthService, error) {
	config, ok := configInstance.(*RequiredHeaderPluginConfig)
	if !ok {
		return nil, errors.New(fmt.Sprintf("unexpected config type %T", configInstance))
	}
	return &RequiredHeaderClient{RequiredHeader: config.RequiredHeader}, nil
}

type RequiredHeaderClient struct {
	RequiredHeader string
}

func (c *RequiredHeaderClient) Start(context.Context) error {
	// no-op
	return nil
}

func (c *RequiredHeaderClient) Authorize(ctx context.Context, request *v2.CheckRequest) (*api.AuthorizationResponse, error) {
	for key, value := range request.Attributes.Request.Http.Headers {
		if key == c.RequiredHeader {
			contextutils.LoggerFrom(ctx).Infow("found required header", "header", key, "value", value)
			return api.AuthorizedResponse(), nil
		}
	}
	contextutils.LoggerFrom(ctx).Infow("required header not found, denying access")
	return api.UnauthorizedResponse(), nil
}
