FROM golang:1.18-alpine as builder

ARG CI_COMMIT_TAG
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/

RUN set -ex \
    CGO_ENABLED=0 go build -o release/goapp \
    -trimpath \
    -ldflags "-w -s \
    -X megpoid.xyz/go/go-skel/internal/version.Tag=${CI_COMMIT_TAG}"

FROM alpine:3.15
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /src/release/app /usr/local/bin/goapp

ENTRYPOINT ["/usr/local/bin/goapp"]
