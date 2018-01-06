#docker run -it --rm -p=8080:8080 foo
#docker build -t foo .
FROM golang
RUN mkdir /app 
ADD go-pdf-bot /app/ 
WORKDIR /app 
ENTRYPOINT ./go-pdf-bot