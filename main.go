package main

import (
	"context"
	"encoding/json"
	"flag"

	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/endpoint"
)

func main() {
	port := flag.String("p", "8080", "Port")
	headlessURL := flag.String("-c", "http://localhost:9222/json", "Chrome headless url")
	queueSize := flag.Int("-s", 100, "Queue size")
	flag.Parse()

	jobQueue := make(chan job, *queueSize)

	worker := newWorker(*headlessURL, jobQueue)
	worker.start()

	//logger
	logger := log.NewLogfmtLogger(os.Stderr)

	//instru
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "default",
		Subsystem: "go_pdf_bot",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "default",
		Subsystem: "go_pdf_bot",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "default",
		Subsystem: "go_pdf_bot",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	//router
	var router Router
	router = jsonrpcRouter{
		map[string]controller{
			"job/create": jobCreateController{worker},
		},
	}

	router = loggingMiddleware{logger, router}
	router = instrumentingMiddleware{requestCount, requestLatency, countResult, router}

	handler := httptransport.NewServer(
		rpcEndpoint(router),
		decodeRequest,
		encodeResponse,
	)
	http.Handle("/rpc", handler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log(http.ListenAndServe(":"+*port, nil))
}

// JsonRpcRuter router

type jsonrpcRequest struct {
	Jsonrpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	ID      string            `json:"id"`
}

type jsonrpcErrorResponse struct {
	Code    int               `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Data    map[string]string `json:"data,omitempty"`
}

type jsonrpcResponse struct {
	Jsonrpc string               `json:"jsonrpc"`
	Result  map[string]string    `json:"result"`
	Error   jsonrpcErrorResponse `json:"error"`
	ID      string               `json:"id"`
}

//Router interface
type Router interface {
	route(ctx context.Context, request jsonrpcRequest) (jsonrpcResponse, error)
}

type jsonrpcRouter struct {
	routes map[string]controller
}

type controller interface {
	action(params map[string]string) map[string]string
}

func (j jsonrpcRouter) route(ctx context.Context, request jsonrpcRequest) (jsonrpcResponse, error) {
	response := jsonrpcResponse{
		"2.0",
		map[string]string{},
		jsonrpcErrorResponse{},
		request.ID,
	}

	if j.routes[request.Method] == nil {
		response.Error.Code = -32601
		response.Error.Message = "Method not found"
		return response, nil
	}

	response.Result = j.routes[request.Method].action(request.Params)
	return response, nil
}

func rpcEndpoint(router Router) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(jsonrpcRequest)
		response, error := router.route(ctx, req)
		return response, error
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req jsonrpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	resp := response.(jsonrpcResponse)
	return json.NewEncoder(w).Encode(resp)
}
