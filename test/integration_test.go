package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadlessRun(t *testing.T) {
	r := send()
	assert.Equal(t, "200 OK", r.Status, "the should be equal")
	assert.Equal(t, `{"jsonrpc":"2.0","result":{},"error":{},"id":"128612876124812"}`, r.Body, "the should be equal")
	assert.Equal(t, []string{"application/json"}, r.Headers["Content-Type"], "the should be equal")
}

func send() HTTPResponse {
	body := `{
		"jsonrpc" : "2.0",
		"method":"job/create",
		"params":{
			"html" :"Hello <b>World</b><i>!</i>"
		},
		"id":"128612876124812"
	}`
	h := &HTTPRequest{}
	return h.Method("Post").URL("http://localhost:8080/rpc").Body(body).Header("Content-Type", "application/json").Send()
}
