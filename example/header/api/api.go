package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/grpc"
)

type RequiredHeaderPlugin struct {
	RequiredHeader string
}

func (p *RequiredHeaderPlugin) NewConfigInstance(ctx context.Context) interface{} {
	return &RequiredHeaderPlugin{}
}

func (p *RequiredHeaderPlugin) GetAuthClient(ctx context.Context, configInstance interface{}) (api.AuthClient, error) {
	config, ok := configInstance.(*RequiredHeaderPlugin)
	if !ok {
		return nil, errors.New(fmt.Sprintf("unexpected config type %T", configInstance))
	}
	return &RequiredHeaderClient{RequiredHeader: config.RequiredHeader}, nil
}

type RequiredHeaderClient struct {
	RequiredHeader string
}

func (c *RequiredHeaderClient) Start() {
	_ = grpc.NewServer()
	// no-op
}

func (c *RequiredHeaderClient) Authorize(ctx context.Context, request *api.Request) (*api.AuthorizationResponse, error) {
	//for key, value := range request.Attributes.Request.Http.Headers {
	//	if key == c.RequiredHeader {
	//		contextutils.LoggerFrom(ctx).Infow("found required header", "header", key, "value", value)
	//		return api.AuthorizedResponse(), nil
	//	}
	//}
	//contextutils.LoggerFrom(ctx).Infow("required header not found, denying access")
	contextutils.LoggerFrom(ctx).Infow("allowing all")
	return api.AuthorizedResponse(), nil
}
