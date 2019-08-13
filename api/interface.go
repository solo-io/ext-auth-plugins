package api

import (
	"context"

	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/googleapis/google/rpc"
)

type StartFunc func(ctx context.Context) error
type AuthorizeFunc func(ctx context.Context, request *envoyauthv2.CheckRequest) (*AuthorizationResponse, error)

// Response returned by authorization services to the Gloo ext-auth server
type AuthorizationResponse struct {
	// Additional user information
	UserInfo UserInfo
	// The result of the authorization process that will be sent back to Envoy
	CheckResponse envoyauthv2.CheckResponse
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
	Authorize(ctx context.Context, request *envoyauthv2.CheckRequest) (*AuthorizationResponse, error)
}

// External authorization plugins must implement this interface
type ExtAuthPlugin interface {
	// Gloo will deserialize the external authorization plugin configuration defined on your Virtual Hosts into the
	// type returned by this function. The returned type MUST be a pointer.
	//
	// For example, given the following plugin configuration:
	//
	//   apiVersion: gateway.solo.io/v1
	//   kind: VirtualService
	//   metadata:
	//     name: test-auth
	//     namespace: gloo-system
	//   spec:
	//     virtualHost:
	//       domains: [...]
	//       routes: [...]
	//       virtualHostPlugins:
	//         extensions:
	//           configs:
	//             extauth:
	//               plugin_auth:
	//                 name: MyAuthPlugin
	//                 config:
	//                   some_key: value-1
	//                   some_struct:
	//                     another_key: value-2
	//
	// the `NewConfigInstance` function on your `ExtAuthPlugin` implementation should return a pointer to
	// the following Go struct:
	//
	//   type MyPluginConfig struct {
	//     SomeKey string
	//	   SomeStruct NestedConfig
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
		CheckResponse: envoyauthv2.CheckResponse{},
	}
}

// Minimal DENIED response
func UnauthorizedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		CheckResponse: envoyauthv2.CheckResponse{
			Status: &rpc.Status{
				Code: int32(rpc.PERMISSION_DENIED),
			},
		},
	}
}

func InternalServerErrorResponse() *AuthorizationResponse {
	resp := UnauthorizedResponse()
	resp.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_DeniedResponse{
		DeniedResponse: &envoyauthv2.DeniedHttpResponse{
			Status: &envoytype.HttpStatus{
				Code: envoytype.StatusCode_InternalServerError,
			},
		},
	}
	return resp
}
