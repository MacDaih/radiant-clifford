FROM golang:1.17.3-alpine

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=arm64 go build -o main .

CMD ["/go/src/app/main"]