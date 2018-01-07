package main

import (
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		pdfBody, _ := ioutil.ReadAll(r.Body)
		err := ioutil.WriteFile("test.pdf", pdfBody, 0644)
		if nil != err {
			panic(err)
		}
	})

	err := http.ListenAndServe(":8833", nil)
	if err != nil {
		panic(err)
	}
}

/**

curl -X POST \
  http://localhost:8080/rpc \
  -H 'content-type: application/json' \
  -d '{
	"jsonrpc" : "2.0",
	"method":"job/create",
	"params":{
		"html" :"Hello <b>World</b>!!",
		"webhook" :"http://localhost:8833/webhook"
	},
	"id":"128612876124812"
}'
**/
