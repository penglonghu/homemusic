# 构建阶段
FROM golang:1.25-alpine AS builder
WORKDIR /app

# 代理配置
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=1
ENV GOOS=linux

# 安装依赖
RUN apk add --no-cache git gcc musl-dev zlib-dev

# 🔥 核心修复：仅复制依赖文件（缓存层！永远不会重复下载）
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 🔥 最后复制源码（源码修改，只编译源码，不重新下依赖）
COPY backend/ ./

# 编译
RUN go build -o homemusic-api main.go

# 运行阶段
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

RUN mkdir -p /app/data
COPY --from=builder /app/homemusic-api .
COPY --from=builder /app/config/config.yaml ./config/

RUN chmod +x /app/homemusic-api
EXPOSE 8080
CMD ["./homemusic-api"]