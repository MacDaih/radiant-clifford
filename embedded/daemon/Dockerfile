FROM golang:1.16.3-alpine3.12

RUN mkdir /daemon

ADD . /daemon

WORKDIR /daemon

RUN GOOS=linux GOARCH=arm64 go build -o main .

CMD ["/daemon/main"]