# 阶段 1: 前端构建
FROM node:14 as frontend-build
WORKDIR /app
# 复制前端代码并安装依赖
COPY mon-app/package*.json ./mon-app/
WORKDIR /app/mon-app
RUN npm install
COPY mon-app/ ./
# 构建静态文件
RUN npm run build

# 阶段 2: 后端构建
FROM golang:1.16 as backend-build
WORKDIR /app
# 复制后端代码并安装依赖
COPY serveur/go.mod serveur/go.sum ./serveur/
WORKDIR /app/serveur
RUN go mod download
COPY serveur/ ./
# 编译 Go 应用
RUN go build -o server
RUN ls -la /app/serveur
RUN ls -la /app/serveur/


# 阶段 3: 创建最终镜像
FROM golang:alpine
WORKDIR /app
# 从前端构建阶段复制静态文件
COPY --from=frontend-build /app/mon-app/build /app/mon-app/build
# 从后端构建阶段复制编译好的二进制文件
COPY --from=backend-build /app/serveur/server /app/serveur/server
# 设置运行命令
CMD ["/app/serveur/server"]s

