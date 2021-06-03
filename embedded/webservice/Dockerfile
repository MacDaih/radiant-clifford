FROM golang:1.16.3-alpine3.12

RUN mkdir /ws

ADD . /ws

WORKDIR /ws

RUN GOOS=linux GOARCH=arm64 go build -o main .

CMD ["/ws/main"]

EXPOSE 8080