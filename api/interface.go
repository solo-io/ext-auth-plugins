package api

import (
	"context"

	pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	googlerpc "github.com/gogo/googleapis/google/rpc"
)

type UserInfo struct {
	UserID string
}

type AuthorizationResponse struct {
	UserInfo      UserInfo
	CheckResponse pb.CheckResponse
}

type AuthClient interface {
	Start()
	Authorize(ctx context.Context, request *pb.CheckRequest) (*AuthorizationResponse, error)
}

type ExtauthPlugin interface {
	NewConfigInstance(ctx context.Context) interface{} //proto message
	GetAuthClient(ctx context.Context, configInstance interface{}) (AuthClient, error)
}

func AuthorizedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		CheckResponse: pb.CheckResponse{},
	}
}

func UnauthorizedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		CheckResponse: pb.CheckResponse{
			Status: &googlerpc.Status{
				Code: int32(googlerpc.PERMISSION_DENIED),
			},
		},
	}
}