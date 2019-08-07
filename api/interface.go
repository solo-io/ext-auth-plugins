package api

import (
	"context"
)

type Request struct {
	Id string
}

type Response struct {
	Authorized bool
}

type UserInfo struct {
	UserID string
}

type AuthorizationResponse struct {
	UserInfo UserInfo
	//CheckResponse pb.CheckResponse
	Response Response
}

// TODO: add ctx to start
type AuthClient interface {
	Start()
	//Authorize(ctx context.Context, request *pb.CheckRequest) (*AuthorizationResponse, error)
	Authorize(ctx context.Context, request *Request) (*AuthorizationResponse, error)
}

type ExtauthPlugin interface {
	NewConfigInstance(ctx context.Context) interface{} //proto message
	GetAuthClient(ctx context.Context, configInstance interface{}) (AuthClient, error)
}

func AuthorizedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		//CheckResponse: pb.CheckResponse{},
		Response: Response{
			Authorized: true,
		},
	}
}

func UnauthorizedResponse() *AuthorizationResponse {
	return &AuthorizationResponse{
		//CheckResponse: pb.CheckResponse{
		//	Status: &googlerpc.Status{
		//		Code: int32(googlerpc.PERMISSION_DENIED),
		//	},
		//},
		Response: Response{
			Authorized: false,
		},
	}
}

func InternalServerErrorResponse() *AuthorizationResponse {
	resp := UnauthorizedResponse()
	//resp.CheckResponse.HttpResponse = &pb.CheckResponse_DeniedResponse{
	//	DeniedResponse: &pb.DeniedHttpResponse{
	//		Status: &envoytype.HttpStatus{
	//			Code: envoytype.StatusCode_InternalServerError,
	//		},
	//	},
	//}

	return resp
}
