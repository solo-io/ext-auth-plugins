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

// External authorization plugins must implement this interface
type ExtAuthPlugin interface {
	// Gloo will deserialize the external authorization plugin configuration defined on your AuthConfig into the
	// type returned by this function. The returned type MUST be a pointer.
	//
	// For example, given the following plugin configuration:
	//
	//  apiVersion: enterprise.gloo.solo.io/v1
	//  kind: AuthConfig
	//  metadata:
	//    name: plugin-auth
	//    namespace: gloo-system
	//  spec:
	//    configs:
	//    - pluginAuth:
	//        name: MyAuthPlugin
	//        config:
	//          someKey: value-1
	//          someStruct:
	//            anotherKey: value-2
	//
	// the `NewConfigInstance` function on your `ExtAuthPlugin` implementation should return a pointer to
	// the following Go struct:
	//
	//   type MyPluginConfig struct {
	//     SomeKey string
	//     SomeStruct NestedConfig
	//   }
	//
	// where `NestedConfig` is:
	//
	//   type NestedConfig struct {
	//     AnotherKey string
	//   }
	//
	// When Gloo invokes this fun function during plugin initialization, it will pass in a context. The context will
	// be cancelled whenever Gloo detects a change in the overall auth configuration and consequently re-initializes
	// all the auth plugins.
	NewConfigInstance(ctx context.Context) (configInstance interface{}, err error)

	// Returns an authorization service instance.
	// The input context is the same as the one passed to `NewConfigInstance` and the same considerations apply.
	// The input configInstance is the same one returned by the `NewConfigInstance` after the plugin configuration has
	// been deserialized into it.
	GetAuthService(ctx context.Context, configInstance interface{}) (AuthService, error)
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
