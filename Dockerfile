# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS build
LABEL "team"="c8s"
LABEL version="1.0"

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.io",direct

WORKDIR /app
ADD . .

RUN go mod tidy
RUN go mod vendor

RUN apk add dos2unix
RUN find . -type f -print0 | xargs -0 dos2unix
RUN go build  -ldflags "-s -w" -o ./build/gpuserver cmd/gpujobserver/gpujobserver.go

ENTRYPOINT [ "./build/gpuserver" ]

# TODO: Smaller image:)
#-ldflags "-s -w"
# FROM scratch As prod
#
# WORKDIR /bin/
#
# COPY --from=0 /app/build/gpuserver .
# #COPY --from=build  /app/build/gpuserver /
# ENTRYPOINT [ "./gpuserver" ]
