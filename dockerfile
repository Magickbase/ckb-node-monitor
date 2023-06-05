# 构建编译环境
FROM golang:1.18 AS build-env
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o monitor

# 运行环境
FROM alpine:latest
RUN apk update && apk add ca-certificates
COPY --from=build-env /app/monitor /app/monitor
WORKDIR /app
CMD ["./monitor"]
