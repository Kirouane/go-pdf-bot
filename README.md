# Go PDF Bot

Generate PDF files from HTML pages **asynchronously**. This tool has been implemented with Golang and inspired by the excelent NodeJS tool [pdf-bot](https://github.com/esbenp/pdf-bot).

## Install

```bash
go get -d ./...
```

## Run

```bash
go build .
go ./go-pdf-bot
```

Go PDF Bot opened a port on 8080 (default port) and it is waiting for HTTP requests like :

```bash
curl -X POST \
  http://localhost:8080/rpc \
  -H 'content-type: application/json' \
  -d '{
    "jsonrpc" : "2.0",
    "method":"job/create",
    "params":{
        "html" :"Hello <b>World</b>!!"
    },
    "id":"128612876124812"
}'
```

You can see the parameter "html" with value `Hello <b>World</b>!!`. Go PDF Bot will generate a PDF file in **storage/pdf** directory based on this html code.

## Webhook

You can tell to Go PDF Bot to post the PDF content on your server through a webhook : 

```bash
curl -X POST \
  http://localhost:8080/rpc \
  -H 'content-type: application/json' \
  -d '{
    "jsonrpc" : "2.0",
    "method":"job/create",
    "params":{
        "html" :"Hello <b>World</b>!!",
        "webhook" :"https://your.website.com"
    },
    "id":"128612876124812"
}'
```

## Example

* Run a Go PDF bot instance.

```bash
go build .
go ./go-pdf-bot
```

* Run a server to collect PDF contents listening on the port 8833.

```bash
go run webhook.go
```

* Make a curl request to generate a pdf and post the coontent on this.

```bash
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
```
