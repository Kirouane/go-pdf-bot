package integration

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

//HTTPRequest request object
type HTTPRequest struct {
	method  string
	url     string
	body    string
	headers map[string]string
}

//Method func
func (r *HTTPRequest) Method(method string) *HTTPRequest {
	r.method = method
	return r
}

//URL func
func (r *HTTPRequest) URL(url string) *HTTPRequest {
	r.url = url
	return r
}

//Body func
func (r *HTTPRequest) Body(body string) *HTTPRequest {
	r.body = body
	return r
}

//Header func
func (r *HTTPRequest) Header(name string, value string) *HTTPRequest {
	if nil == r.headers {
		r.headers = map[string]string{}
	}
	r.headers[name] = value
	return r
}

//Send func
func (r *HTTPRequest) Send() HTTPResponse {
	req, err := http.NewRequest(r.method, r.url, bytes.NewBuffer([]byte(r.body)))

	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return HTTPResponse{
		resp.Status,
		strings.TrimSpace(string(body)),
		map[string][]string(resp.Header),
	}
}

//HTTPResponse request object
type HTTPResponse struct {
	Status  string
	Body    string
	Headers map[string][]string
}
