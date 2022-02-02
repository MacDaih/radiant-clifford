FROM golang:1.16.3-alpine3.12

WORKDIR /go/src/ws

COPY . .

RUN go mod download
RUN GOOS=linux GOARCH=arm64 go build -o main .

EXPOSE 8080

CMD ["/go/src/ws/main"]