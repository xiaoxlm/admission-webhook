FROM golang:1.20-alpine
WORKDIR /go/src
COPY ./ ./
ARG GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 go  build -a -ldflags '-s -w' -o webhook
EXPOSE 80
ENTRYPOINT ["./webhook"]