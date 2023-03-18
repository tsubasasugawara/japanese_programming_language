FROM golang:alpine

RUN apk update && apk add vim git curl bash

ARG USERNAME=user
ARG GROUPNAME=user
ARG UID=1000
ARG GID=1000
RUN addgroup -g $GID -S $GROUPNAME
RUN adduser -u $UID -S $USERNAME -G $GROUPNAME
USER $USERNAME
WORKDIR /go/src

COPY src .

CMD ["bash"]
