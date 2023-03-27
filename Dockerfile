FROM golang:alpine

RUN apk update && apk add vim curl bash

WORKDIR /go/src
COPY src .

CMD ["bash"]
