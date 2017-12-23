#docker run -it --rm -p=8080:8080 foo
#docker build -t foo .
FROM golang
WORKDIR .
COPY . .
RUN go get -d ./...
EXPOSE 8080

ENTRYPOINT go run $(ls -1 *.go | grep -v _test.go)
