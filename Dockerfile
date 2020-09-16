FROM golang:1.15-alpine as builder

ARG CI_TAG
ARG BUILD_NUMBER
ARG BUILD_COMMIT_SHORT
ARG CI_BUILD_CREATED
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
ENV CGO_ENABLED 0
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/
RUN go build -o release/app \
   -ldflags "-w -s \
   -X main.Version=${CI_TAG} \
   -X main.BuildNumber=${BUILD_NUMBER} \
   -X main.Commit=${BUILD_COMMIT_SHORT} \
   -X main.BuildTime=${CI_BUILD_CREATED}"

FROM alpine:3.12
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /src/release/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
