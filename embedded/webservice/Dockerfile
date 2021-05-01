FROM golang:1.16.3-alpine3.12

RUN mkdir /ws

ADD . /ws

WORKDIR /ws

RUN go build -o main .

CMD ["/ws/main"]