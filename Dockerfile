# 使用官方Golang镜像作为构建环境
FROM golang:1.23 AS builder

# 设置 Go 环境变量
ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=off 
    
    
# 设置工作目录
WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./ 

# 下载依赖
RUN go mod tidy && \  
    go mod download

# 复制项目文件
COPY . .

# 编译项目
RUN CGO_ENABLED=0 GOOS=linux go build -o vireo-gin-admin

# 使用轻量级Alpine镜像作为运行时环境
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/vireo-gin-admin .

# 复制配置文件
COPY config/config.yaml ./config/

# 暴露服务端口
EXPOSE 8080

# 启动命令
CMD ["./vireo-gin-admin"]