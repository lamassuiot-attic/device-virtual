package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return &instrumentingMiddleware{
			requestCount:   counter,
			requestLatency: latency,
			next:           next,
		}
	}
}

func (mw *instrumentingMiddleware) PostSendMessage(ctx context.Context, message string, topic string) (err error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "PostSendMessage").Add(1)
		mw.requestLatency.With("method", "PostSendMessage").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.PostSendMessage(ctx, message, topic)
}

func (mw *instrumentingMiddleware) PostConnect(ctx context.Context, authKey string, authCRT string, brokerURL string, clientID string) (err error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "PostConnect").Add(1)
		mw.requestLatency.With("method", "PostConnect").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.PostConnect(ctx, authKey, authCRT, brokerURL, clientID)
}

func (mw *instrumentingMiddleware) PostDisconnect(ctx context.Context) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "PostDisconnect").Add(1)
		mw.requestLatency.With("method", "PostDisconnect").Observe(time.Since(begin).Seconds())
	}(time.Now())

	mw.next.PostDisconnect(ctx)
}
