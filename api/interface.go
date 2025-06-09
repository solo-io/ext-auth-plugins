package api

import (
	"context"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type StartFunc func(ctx context.Context) error
type AuthorizeFunc func(ctx context.Context, request *AuthorizationRequest) (*AuthorizationResponse, error)

// Response returned by authorization services to the Gloo ext-auth server
type AuthorizationResponse struct {
	// Additional user information
	UserInfo UserInfo
	// Additional user information
	ApiProductInfo ApiProductInfo
	// The result of the authorization process that will be sent back to Envoy
	CheckResponse envoy_service_auth_v3.CheckResponse
}

type AuthorizationRequest struct {
	// The request that needs to be authorized
	CheckRequest *envoy_service_auth_v3.CheckRequest
	State        map[string]interface{}
}

func (a *AuthorizationRequest) SetState(key string, value interface{}) {
	if a.State == nil {
		a.State = make(map[string]interface{})
	}
	a.State[key] = value
}

func (a *AuthorizationRequest) GetState(key string) interface{} {
	if a.State == nil {
		return nil
	}
	return a.State[key]
}

// TODO(marco): consider moving this to the ext-auth-service
// Can be used to set an additional header on authorized requests.
type UserInfo struct {
	UserID string
}

type ApiProductInfo struct {
	UsagePlan string
}

// AuthService instances are responsible for authorizing individual requests.
type AuthService interface {
	// This function will be called when the authorization service is started. Whenever the auth configuration changes,
	// Gloo will start a new instance of the AuthService and signal the termination of the previous one by cancelling
	// the provided context.
	Start(ctx context.Context) error

	// Each time a request hits the auth service, this function will be invoked to decide whether nor not to authorize it.
	// If a non-nil error is returned, the request will be denied.
	Authorize(ctx context.Context, request *AuthorizationRequest) (*AuthorizationResponse, error)
}

// Minimal OK response
func AuthorizedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		CheckResponse: envoy_service_auth_v3.CheckResponse{
			Status: &status.Status{
				Code: int32(codes.OK),
			},
		},
	}
}

// Minimal FORBIDDEN (403) response
func UnauthorizedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		CheckResponse: envoy_service_auth_v3.CheckResponse{
			Status: &status.Status{
				Code: int32(codes.PermissionDenied),
			},
		},
	}
}

// Minimal UNAUTHORIZED (401) response
func UnauthenticatedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		CheckResponse: envoy_service_auth_v3.CheckResponse{
			Status: &status.Status{
				Code: int32(codes.Unauthenticated),
			},
			HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
				DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
					Status: &envoy_type_v3.HttpStatus{
						Code: envoy_type_v3.StatusCode_Unauthorized,
					},
				},
			},
		},
	}
}

func InternalServerErrorResponse() *AuthorizationResponse {
	resp := UnauthorizedResponse()
	resp.CheckResponse.HttpResponse = &envoy_service_auth_v3.CheckResponse_DeniedResponse{
		DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
			Status: &envoy_type_v3.HttpStatus{
				Code: envoy_type_v3.StatusCode_InternalServerError,
			},
		},
	}
	return resp
}
