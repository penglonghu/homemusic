# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY frontend/package.json ./
# 加速npm安装（国内镜像+跳过无用步骤）
RUN npm config set registry https://registry.npmmirror.com/
RUN npm install --no-fund --no-audit --verbose

# 复制源码
COPY frontend/ ./
RUN npm run build

# 运行阶段
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]