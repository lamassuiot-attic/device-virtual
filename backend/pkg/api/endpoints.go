package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PostSendMessage endpoint.Endpoint
	PostConnect     endpoint.Endpoint
	PostDisconnect  endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostConnect:     MakePostConnect(s),
		PostDisconnect:  MakePostDisconnect(s),
		PostSendMessage: MakePostSendMessage(s),
	}
}

func MakePostConnect(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postConnectRequest)
		err = s.PostConnect(ctx, req.AuthKey, req.AuthCRT, req.BrokerURL, req.ClientID)
		return postConnectResponse{Err: err}, nil
	}
}

func MakePostDisconnect(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(postDisconnectRequest)
		s.PostDisconnect(ctx)
		return postDisconnectResponse{}, nil
	}
}

func MakePostSendMessage(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postSendMessageRequest)
		err = s.PostSendMessage(ctx, req.Message)
		return postSendMessageResponse{Err: err}, nil
	}
}

type postConnectRequest struct {
	AuthKey   string `json:"authKey"`
	AuthCRT   string `json:"authCRT"`
	BrokerURL string `json:"brokerURL"`
	ClientID  string `json:"clientID"`
}

type postConnectResponse struct {
	Err error `json:"error"`
}

func (r postConnectResponse) error() error { return r.Err }

type postDisconnectRequest struct{}

type postDisconnectResponse struct{}

type postSendMessageRequest struct {
	Message string `json:"message"`
}

type postSendMessageResponse struct {
	Err error `json:"error"`
}

func (r postSendMessageResponse) error() error { return r.Err }
