FROM golang:1.18-alpine as builder

ARG CI_COMMIT_TAG
ARG CI_COMMIT_BRANCH
ARG CI_COMMIT_SHA
ARG CI_PIPELINE_CREATED_AT
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
    -X main.Version=${CI_COMMIT_TAG:-$CI_COMMIT_BRANCH} \
    -X main.Commit=$(echo "$CI_COMMIT_SHA" | cut -c1-8) \
    -X main.BuildTime=${CI_PIPELINE_CREATED_AT}"

FROM alpine:3.15
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /src/release/app /usr/local/bin/goapp

ENTRYPOINT ["/usr/local/bin/goapp"]
