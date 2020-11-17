FROM golang:1.15-alpine

ARG GO_WORK_DIR
ARG GOPROXY

EXPOSE 8000

ENV TZ Europe/Paris

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GOPATH /go
ENV GOCACHE /tmp/go-cache
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin"

ENV GO111MODULE on

COPY go.mod go.sum ./
COPY docker/dev/root/* /.
RUN chmod +x /entrypoint.sh

RUN \ 
    apk add --no-cache bash git openssh

CMD ["go", "run", "main.go"]