#!/bin/bash

# 博客系统重新部署脚本
# 用途：停止所有服务，清理容器和镜像，重新构建并启动

set -e  # 遇到错误立即退出

echo "======================================"
echo "  博客系统重新部署脚本"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker是否运行
check_docker() {
    log_info "检查Docker是否运行..."
    if ! docker info > /dev/null 2>&1; then
        log_error "Docker未运行，请先启动Docker"
        exit 1
    fi
    log_info "Docker运行正常"
}

# 停止并删除所有相关容器
stop_containers() {
    log_info "停止并删除所有容器..."
    
    if [ "$(docker ps -aq -f name=blog-)" ]; then
        docker stop $(docker ps -aq -f name=blog-) 2>/dev/null || true
        docker rm $(docker ps -aq -f name=blog-) 2>/dev/null || true
        log_info "容器已停止并删除"
    else
        log_warn "没有找到运行中的容器"
    fi
}

# 删除相关镜像（可选）
remove_images() {
    log_info "删除旧的镜像..."
    
    docker rmi blog-backend-dev 2>/dev/null || true
    docker rmi blog-frontend-dev 2>/dev/null || true
    
    log_info "镜像已删除"
}

# 清理dangling镜像和构建缓存
clean_docker() {
    log_info "清理Docker dangling资源..."
    
    docker image prune -f > /dev/null 2>&1 || true
    
    log_info "清理完成"
}

# 生成Swagger文档
generate_swagger() {
    log_info "生成Swagger文档..."
    
    cd backend
    
    # 检查swag是否安装
    if ! command -v swag &> /dev/null; then
        log_warn "swag未安装，尝试安装..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    
    # 生成文档
    if ! swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal 2>&1; then
        log_warn "Swagger文档生成失败，将在容器内生成"
    else
        log_info "Swagger文档生成成功"
    fi
    
    cd ..
}

# 构建并启动服务
start_services() {
    log_info "构建并启动服务（使用缓存）..."
    
    # 使用docker-compose.dev.yml启动开发环境
    docker compose -f docker-compose.dev.yml build --no-cache
    docker compose -f docker-compose.dev.yml up -d
    
    log_info "服务启动中..."
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务启动..."
    
    # 等待后端服务
    MAX_RETRIES=30
    RETRY_COUNT=0
    
    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            log_info "后端服务已就绪"
            break
        fi
        
        RETRY_COUNT=$((RETRY_COUNT+1))
        echo -n "."
        sleep 2
    done
    
    echo ""
    
    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        log_error "后端服务启动超时"
        return 1
    fi
    
    # 等待前端服务
    sleep 3
    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        log_info "前端服务已就绪"
    else
        log_warn "前端服务可能还在启动中"
    fi
}

# 显示服务状态
show_status() {
    log_info "服务状态："
    echo ""
    docker compose -f docker-compose.dev.yml ps
    echo ""
}

# 显示访问信息
show_info() {
    echo ""
    echo "======================================"
    echo -e "${GREEN}  部署完成！${NC}"
    echo "======================================"
    echo ""
    echo "访问地址："
    echo "  前端：http://localhost:5173"
    echo "  后端API：http://localhost:8080"
    echo "  Swagger文档：http://localhost:8080/swagger/index.html"
    echo ""
    echo "默认管理员账号："
    echo "  用户名：admin"
    echo "  密码：admin@123"
    echo ""
    echo "查看日志："
    echo "  后端：docker logs -f blog-backend-dev"
    echo "  前端：docker logs -f blog-frontend-dev"
    echo ""
    echo "停止服务："
    echo "  docker compose -f docker-compose.dev.yml down"
    echo "======================================"
}

# 主函数
main() {
    # 切换到项目根目录
    cd "$(dirname "$0")"
    
    # 执行部署步骤
    check_docker
    stop_containers
    remove_images
    clean_docker
    generate_swagger
    start_services
    wait_for_services
    show_status
    show_info
}

# 运行主函数
main

