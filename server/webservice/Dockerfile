FROM golang:1.17.3-alpine

WORKDIR /go/src/ws

COPY . .

RUN go mod download
RUN go build -o main ./cmd

CMD ["/go/src/ws/main"]