package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r.Methods("POST").Path("/v1/device/connect").Handler(httptransport.NewServer(
		e.PostConnect,
		decodePostConnectRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/device/disconnect").Handler(httptransport.NewServer(
		e.PostDisconnect,
		decodePostDisconnectRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/device/message").Handler(httptransport.NewServer(
		e.PostSendMessage,
		decodePostSendMessageRequest,
		encodeResponse,
		options...,
	))
	return r
}

type errorer interface {
	error() error
}

func decodePostSendMessageRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var reqData postSendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return nil, err
	}
	return reqData, nil
}

func decodePostConnectRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var reqData postConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return nil, err
	}
	return reqData, nil
}

func decodePostDisconnectRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var reqData postDisconnectRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return nil, err
	}
	return reqData, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)

		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	http.Error(w, err.Error(), codeFrom(err))
}

func codeFrom(err error) int {
	switch err {
	case ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
