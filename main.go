package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

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

type jsonrpcRouter struct {
	routes map[string]controller
}

type controller interface {
	action(params map[string]string) map[string]string
}

func (j *jsonrpcRouter) route(_ context.Context, request jsonrpcRequest) (jsonrpcResponse, error) {
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

var jobQueue = make(chan job, 10)

func main() {

	worker := newWorker(jobQueue)
	worker.start()

	router := jsonrpcRouter{
		map[string]controller{},
	}
	router.routes["job/create"] = jobCreateController{worker}

	handler := httptransport.NewServer(
		rpcEndpoint(router),
		decodeRequest,
		encodeResponse,
	)
	http.Handle("/rpc", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rpcEndpoint(router jsonrpcRouter) endpoint.Endpoint {
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
