#!/bin/bash

# 一键部署脚本
# 自动检查环境、构建镜像、启动服务

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 博客系统一键部署脚本${NC}"
echo ""

# 检查Docker
check_docker() {
    echo -e "${YELLOW}🔍 检查Docker环境...${NC}"
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker未安装，请先安装Docker${NC}"
        exit 1
    fi
    if ! docker info &> /dev/null; then
        echo -e "${RED}❌ Docker服务未运行，请启动Docker服务${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker环境正常${NC}"
}

# 检查Docker Compose
check_docker_compose() {
    echo -e "${YELLOW}🔍 检查Docker Compose...${NC}"
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}❌ Docker Compose未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker Compose正常${NC}"
}

# 检查端口占用
check_ports() {
    echo -e "${YELLOW}🔍 检查端口占用...${NC}"
    PORTS=(80 443 5432 6379 8080)
    OCCUPIED=()
    
    for port in "${PORTS[@]}"; do
        if lsof -Pi :${port} -sTCP:LISTEN -t >/dev/null 2>&1; then
            OCCUPIED+=(${port})
        fi
    done
    
    if [ ${#OCCUPIED[@]} -gt 0 ]; then
        echo -e "${YELLOW}⚠️  以下端口已被占用: ${OCCUPIED[*]}${NC}"
        read -p "继续部署? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        echo -e "${GREEN}✅ 端口检查通过${NC}"
    fi
}

# 检查配置文件
check_config() {
    echo -e "${YELLOW}🔍 检查配置文件...${NC}"
    if [ ! -f "backend/config/config.yaml" ]; then
        if [ -f "backend/config/config.example.yaml" ]; then
            echo -e "${YELLOW}⚠️  配置文件不存在，从示例文件创建...${NC}"
            cp backend/config/config.example.yaml backend/config/config.yaml
            echo -e "${YELLOW}⚠️  请编辑 backend/config/config.yaml 配置数据库等信息${NC}"
            read -p "已创建配置文件，是否继续? (y/N) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        else
            echo -e "${RED}❌ 配置文件不存在且示例文件也不存在${NC}"
            exit 1
        fi
    fi
    echo -e "${GREEN}✅ 配置文件检查通过${NC}"
}

# 创建必要的目录
create_directories() {
    echo -e "${YELLOW}📁 创建必要的目录...${NC}"
    mkdir -p data/postgres
    mkdir -p data/redis
    mkdir -p backend/uploads
    mkdir -p backend/backups
    mkdir -p backend/logs
    mkdir -p docker/nginx/ssl
    mkdir -p docker/nginx/conf.d
    echo -e "${GREEN}✅ 目录创建完成${NC}"
}

# 构建镜像
build_images() {
    echo -e "${YELLOW}🔨 构建Docker镜像...${NC}"
    if command -v docker-compose &> /dev/null; then
        docker-compose build --no-cache
    else
        docker compose build --no-cache
    fi
    echo -e "${GREEN}✅ 镜像构建完成${NC}"
}

# 启动服务
start_services() {
    echo -e "${YELLOW}🚀 启动服务...${NC}"
    if command -v docker-compose &> /dev/null; then
        docker-compose up -d
    else
        docker compose up -d
    fi
    echo -e "${GREEN}✅ 服务启动完成${NC}"
}

# 等待服务就绪
wait_for_services() {
    echo -e "${YELLOW}⏳ 等待服务就绪...${NC}"
    sleep 5
    
    # 检查后端健康状态
    max_attempts=30
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -f http://localhost:8080/health &> /dev/null; then
            echo -e "${GREEN}✅ 后端服务就绪${NC}"
            break
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${YELLOW}⚠️  后端服务启动超时，请检查日志${NC}"
    fi
}

# 显示服务状态
show_status() {
    echo ""
    echo -e "${GREEN}✅ 部署完成!${NC}"
    echo ""
    echo -e "${BLUE}📊 服务状态:${NC}"
    if command -v docker-compose &> /dev/null; then
        docker-compose ps
    else
        docker compose ps
    fi
    echo ""
    echo -e "${BLUE}🌐 访问地址:${NC}"
    echo "  前端: http://localhost"
    echo "  后端API: http://localhost:8080"
    echo "  Swagger文档: http://localhost:8080/swagger/index.html"
    echo ""
    echo -e "${BLUE}📝 常用命令:${NC}"
    echo "  查看日志: docker-compose logs -f"
    echo "  停止服务: docker-compose down"
    echo "  重启服务: docker-compose restart"
    echo ""
}

# 主流程
main() {
    check_docker
    check_docker_compose
    check_ports
    check_config
    create_directories
    
    read -p "是否构建新镜像? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        build_images
    fi
    
    start_services
    wait_for_services
    show_status
}

# 执行主流程
main

