#!/bin/bash

# 健康检查脚本
# 检查所有服务的健康状态

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查结果
ALL_HEALTHY=true

echo -e "${BLUE}🏥 服务健康检查${NC}"
echo ""

# 检查函数
check_service() {
    local name=$1
    local url=$2
    local timeout=${3:-5}
    
    if curl -f -s --max-time ${timeout} ${url} > /dev/null 2>&1; then
        echo -e "${GREEN}✅ ${name}: 健康${NC}"
        return 0
    else
        echo -e "${RED}❌ ${name}: 不健康${NC}"
        ALL_HEALTHY=false
        return 1
    fi
}

# 检查Docker容器
check_container() {
    local name=$1
    
    if docker ps --format '{{.Names}}' | grep -q "^${name}$"; then
        local status=$(docker inspect --format='{{.State.Status}}' ${name} 2>/dev/null)
        if [ "$status" == "running" ]; then
            echo -e "${GREEN}✅ 容器 ${name}: 运行中${NC}"
            return 0
        else
            echo -e "${RED}❌ 容器 ${name}: 状态异常 (${status})${NC}"
            ALL_HEALTHY=false
            return 1
        fi
    else
        echo -e "${RED}❌ 容器 ${name}: 未运行${NC}"
        ALL_HEALTHY=false
        return 1
    fi
}

# 检查数据库连接
check_database() {
    echo -e "${YELLOW}🔍 检查数据库连接...${NC}"
    if docker exec blog-postgres pg_isready -U blog_user -d blog_db > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL: 连接正常${NC}"
    else
        echo -e "${RED}❌ PostgreSQL: 连接失败${NC}"
        ALL_HEALTHY=false
    fi
}

# 检查Redis连接
check_redis() {
    echo -e "${YELLOW}🔍 检查Redis连接...${NC}"
    if docker exec blog-redis redis-cli ping > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Redis: 连接正常${NC}"
    else
        echo -e "${RED}❌ Redis: 连接失败${NC}"
        ALL_HEALTHY=false
    fi
}

# 主检查流程
main() {
    echo -e "${BLUE}📦 检查Docker容器...${NC}"
    check_container "blog-postgres"
    check_container "blog-redis"
    check_container "blog-backend"
    check_container "blog-frontend"
    check_container "blog-nginx"
    
    echo ""
    echo -e "${BLUE}🌐 检查HTTP服务...${NC}"
    check_service "后端API" "http://localhost:8080/health" 5
    check_service "前端页面" "http://localhost/" 5
    check_service "Nginx代理" "http://localhost/health" 5
    
    echo ""
    echo -e "${BLUE}💾 检查数据服务...${NC}"
    check_database
    check_redis
    
    echo ""
    if [ "$ALL_HEALTHY" = true ]; then
        echo -e "${GREEN}✅ 所有服务健康检查通过!${NC}"
        exit 0
    else
        echo -e "${RED}❌ 部分服务健康检查失败，请检查日志${NC}"
        echo ""
        echo -e "${YELLOW}📝 查看日志命令:${NC}"
        echo "  docker-compose logs [service_name]"
        exit 1
    fi
}

# 执行检查
main

