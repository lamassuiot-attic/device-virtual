package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type Middleware func(Service) Service

func LoggingMidleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMidleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMidleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMidleware) PostSendMessage(ctx context.Context, message string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostSendMessage", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostSendMessage(ctx, message)
}

func (mw loggingMidleware) PostConnect(ctx context.Context, authKey string, authCRT string, brokerURL string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostConnect", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostConnect(ctx, authKey, authCRT, brokerURL)
}

func (mw loggingMidleware) PostDisconnect(ctx context.Context) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostDisconnect", "took", time.Since(begin))
	}(time.Now())
	mw.next.PostDisconnect(ctx)
}
