FROM golang:1.20-alpine as build
WORKDIR /go/src
COPY ./ ./
ARG GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 go  build -a -ldflags '-s -w' -o webhook

FROM alpine:3.18.3
COPY --from=build /go/src/webhook /go/bin/webhook
EXPOSE 80
WORKDIR /go/bin
ENTRYPOINT ["/go/bin/webhook"]