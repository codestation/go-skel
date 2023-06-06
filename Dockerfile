FROM golang:1.20-alpine as builder

ARG CI_COMMIT_TAG
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/

RUN set -ex; \
    CGO_ENABLED=0 go build -o release/goapp \
    -trimpath \
    -tags viper_yaml3 \
    -ldflags "-w -s \
    -X megpoid.dev/go/go-skel/version.Tag=${CI_COMMIT_TAG}"

FROM alpine:3.18
LABEL maintainer="codestation <codestation@megpoid.dev>"

RUN apk add --no-cache ca-certificates tzdata

RUN set -eux; \
    addgroup -S runner -g 1000; \
    adduser -S runner -G runner -u 1000

COPY --from=builder /src/release/goapp /usr/local/bin/goapp

USER runner

CMD ["/usr/local/bin/goapp", "serve"]
