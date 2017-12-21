package integration

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeadlessRun(t *testing.T) {
	go webhook(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, []string{"application/pdf"}, r.Header["Content-Type"], "the should be equal")
		assert.NotNil(t, body, "Should be not null")
	})

	r := send()

	assert.Equal(t, "200 OK", r.Status, "the should be equal")
	assert.Equal(t, `{"jsonrpc":"2.0","result":{},"error":{},"id":"128612876124812"}`, r.Body, "the should be equal")
	assert.Equal(t, []string{"application/json"}, r.Headers["Content-Type"], "the should be equal")
	time.Sleep(time.Second * 1)
}

func send() HTTPResponse {
	body := `{
		"jsonrpc" : "2.0",
		"method":"job/create",
		"params":{
			"html" :"Hello <b>World</b><i>!</i>",
			"webhook":"http://localhost:8832/test"
		},
		"id":"128612876124812"
	}`
	h := &HTTPRequest{}
	return h.Method("Post").URL("http://localhost:8080/rpc").Body(body).Header("Content-Type", "application/json").Send()
}

func webhook(t *testing.T, callback func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc("/test", callback)
	err := http.ListenAndServe(":8832", nil)
	if err != nil {
		panic(err)
	}
}
