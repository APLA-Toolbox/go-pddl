
FROM golang:1.15-alpine

ENV TZ Europe/Paris

LABEL maintainer="Erwin Lejeune <erwin.lejeune15@gmail.com>"

RUN apk add --no-cache ca-certificates \
    dpkg \
    gcc \
    git \
    musl-dev \
    bash \
    tzdata


RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GOPATH /go
ENV GOCACHE /tmp/go-cache
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin"

ENV GO111MODULE on

ENV DELVE_VERSION v1.5.0
ENV AIR_VERSION v1.15.1
ENV MOCK_VERSION v1.4.3

RUN go get github.com/golang/mock/gomock@$MOCK_VERSION && \
    go get github.com/golang/mock/mockgen@$MOCK_VERSION && \
    go get github.com/go-delve/delve/cmd/dlv@$DELVE_VERSION && \
    go get github.com/cosmtrek/air@$AIR_VERSION

WORKDIR /go/src/github.com/guilyx/go-pddl
ARG LOG_DIR=/go/src/github.com/guilyx/go-pddl/logs
RUN mkdir -p ${LOG_DIR}
ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log

COPY go.mod ./
COPY go.sum ./
COPY docker/entrypoint.sh /.

RUN chmod +x /entrypoint.sh
RUN chmod +x ./scripts/air.sh

EXPOSE 8080

VOLUME [${LOG_DIR}]

ENTRYPOINT ["/entrypoint.sh"]
