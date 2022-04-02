FROM golang:1.17.3-alpine

WORKDIR /go/src/ws

COPY . .

RUN go mod download
RUN go build -o main .

EXPOSE 8080

CMD ["/go/src/ws/main"]