package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
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
