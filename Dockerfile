FROM golang:alpine as buidler

WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server ./cmd/main.go

FROM alpine:latest

LABEL MAINTAINER="axu9417@gmail.com"

# 添加时区数据
RUN apk add --no-cache tzdata

# 设置时区环境变量
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=0 /app/conf/test/conf.yaml ./conf/test/conf.yaml
COPY --from=0 /app/.kube/config ./.kube/config
COPY --from=0 /app/server ./

EXPOSE 8083
EXPOSE 8082

ENTRYPOINT ./server