package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Router
}

func (mw loggingMiddleware) route(ctx context.Context, request jsonrpcRequest) (response jsonrpcResponse, err error) {
	defer func(begin time.Time) {
		paramsString, _ := json.Marshal(request.Params)
		mw.logger.Log(
			"id", request.ID,
			"method", request.Method,
			"params", string(paramsString),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	response, err = mw.next.route(ctx, request)
	return
}

//instru

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           Router
}

func (mw instrumentingMiddleware) route(ctx context.Context, request jsonrpcRequest) (response jsonrpcResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", request.Method, "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	response, err = mw.next.route(ctx, request)
	return
}
