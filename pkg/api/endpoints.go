package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type Endpoints struct {
	HealthEndpoint  endpoint.Endpoint
	PostSendMessage endpoint.Endpoint
	PostConnect     endpoint.Endpoint
	PostDisconnect  endpoint.Endpoint
}

func MakeServerEndpoints(s Service, otTracer stdopentracing.Tracer) Endpoints {
	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = MakeHealthEndpoint(s)
		healthEndpoint = opentracing.TraceServer(otTracer, "Health")(healthEndpoint)
	}
	var postConnectEndpoint endpoint.Endpoint
	{
		postConnectEndpoint = MakePostConnect(s)
		postConnectEndpoint = opentracing.TraceServer(otTracer, "PostConnect")(postConnectEndpoint)
	}
	var postDisconnectEndpoint endpoint.Endpoint
	{
		postDisconnectEndpoint = MakePostDisconnect(s)
		postDisconnectEndpoint = opentracing.TraceServer(otTracer, "PostDisconnect")(postDisconnectEndpoint)
	}
	var postSendMessageEndpoint endpoint.Endpoint
	{
		postSendMessageEndpoint = MakePostSendMessage(s)
		postSendMessageEndpoint = opentracing.TraceServer(otTracer, "PostSendMessage")(postSendMessageEndpoint)
	}
	return Endpoints{
		HealthEndpoint:  healthEndpoint,
		PostConnect:     postConnectEndpoint,
		PostDisconnect:  postDisconnectEndpoint,
		PostSendMessage: postSendMessageEndpoint,
	}
}

func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		healthy := s.Health(ctx)
		return healthResponse{Healthy: healthy}, nil
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
		err = s.PostSendMessage(ctx, req.Message, req.Topic)
		return postSendMessageResponse{Err: err}, nil
	}
}

type healthRequest struct{}

type healthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
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
	Topic   string `json:"topic"`
}

type postSendMessageResponse struct {
	Err error `json:"error"`
}

func (r postSendMessageResponse) error() error { return r.Err }
