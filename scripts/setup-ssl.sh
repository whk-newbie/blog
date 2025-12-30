#!/bin/bash

# SSL证书申请和管理脚本
# 使用Let's Encrypt自动申请和续期HTTPS证书

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
DOMAIN=""
EMAIL=""
CERTBOT_IMAGE="certbot/certbot:latest"
WEBROOT_PATH="/var/www/certbot"
CERT_PATH="/etc/letsencrypt"

echo -e "${GREEN}🔒 SSL证书管理脚本${NC}"
echo ""

# 检查参数
if [ "$1" == "renew" ]; then
    echo -e "${YELLOW}🔄 续期SSL证书...${NC}"
    docker run --rm \
        -v "$(pwd)/docker/nginx/certbot:/etc/letsencrypt" \
        -v "$(pwd)/docker/nginx/certbot/www:/var/www/certbot" \
        ${CERTBOT_IMAGE} renew --quiet
    echo -e "${GREEN}✅ 证书续期完成${NC}"
    echo -e "${YELLOW}⚠️  请重启Nginx容器以应用新证书: docker-compose restart nginx${NC}"
    exit 0
fi

# 交互式输入
if [ -z "$DOMAIN" ]; then
    read -p "请输入域名 (例如: example.com): " DOMAIN
fi

if [ -z "$EMAIL" ]; then
    read -p "请输入邮箱地址 (用于证书到期提醒): " EMAIL
fi

# 验证输入
if [ -z "$DOMAIN" ] || [ -z "$EMAIL" ]; then
    echo -e "${RED}❌ 域名和邮箱不能为空${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}📋 配置信息:${NC}"
echo "  域名: $DOMAIN"
echo "  邮箱: $EMAIL"
echo ""

read -p "确认信息正确? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}❌ 已取消${NC}"
    exit 1
fi

# 创建必要的目录
echo -e "${YELLOW}📁 创建目录...${NC}"
mkdir -p docker/nginx/certbot/www
mkdir -p docker/nginx/certbot/conf
mkdir -p docker/nginx/certbot/logs

# 检查80端口是否开放
echo -e "${YELLOW}🔍 检查80端口...${NC}"
if ! nc -z localhost 80 2>/dev/null; then
    echo -e "${YELLOW}⚠️  80端口未开放，请确保Nginx服务正在运行${NC}"
    read -p "继续? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 申请证书
echo -e "${YELLOW}📜 申请SSL证书...${NC}"
docker run --rm \
    -v "$(pwd)/docker/nginx/certbot:/etc/letsencrypt" \
    -v "$(pwd)/docker/nginx/certbot/www:/var/www/certbot" \
    ${CERTBOT_IMAGE} certonly \
    --webroot \
    --webroot-path=${WEBROOT_PATH} \
    --email ${EMAIL} \
    --agree-tos \
    --no-eff-email \
    -d ${DOMAIN} \
    -d www.${DOMAIN} || {
    echo -e "${RED}❌ 证书申请失败${NC}"
    exit 1
}

echo -e "${GREEN}✅ 证书申请成功!${NC}"
echo ""

# 更新Nginx配置
echo -e "${YELLOW}⚙️  更新Nginx配置...${NC}"
NGINX_CONF="docker/nginx/nginx.conf"

# 检查证书文件是否存在
if [ ! -f "docker/nginx/certbot/live/${DOMAIN}/fullchain.pem" ]; then
    echo -e "${RED}❌ 证书文件不存在${NC}"
    exit 1
fi

# 创建SSL配置文件
SSL_CONF="docker/nginx/conf.d/ssl-${DOMAIN}.conf"
cat > ${SSL_CONF} <<EOF
# HTTPS服务器配置 - ${DOMAIN}
server {
    listen 443 ssl http2;
    server_name ${DOMAIN} www.${DOMAIN};

    # SSL证书配置（Let's Encrypt）
    ssl_certificate /etc/letsencrypt/live/${DOMAIN}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/${DOMAIN}/privkey.pem;

    # SSL安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # 前端
    location / {
        proxy_pass http://frontend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # 后端API
    location /api {
        proxy_pass http://backend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;

        # WebSocket支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";

        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康检查
    location /health {
        proxy_pass http://backend;
        access_log off;
    }

    # Swagger API文档
    location /swagger {
        proxy_pass http://backend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }

    # 上传文件
    location /uploads {
        proxy_pass http://backend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

# 更新主配置文件，启用HTTP到HTTPS重定向
sed -i 's/# return 301 https:\/\/\$host\$request_uri;/return 301 https:\/\/$host$request_uri;/' ${NGINX_CONF}

echo -e "${GREEN}✅ Nginx配置已更新${NC}"
echo ""

# 更新docker-compose.yml中的Nginx volumes
echo -e "${YELLOW}📝 更新docker-compose.yml...${NC}"
if ! grep -q "nginx_certs:/etc/letsencrypt" docker-compose.yml; then
    echo -e "${YELLOW}⚠️  请手动更新docker-compose.yml中的Nginx volumes配置${NC}"
fi

echo ""
echo -e "${GREEN}✅ SSL证书配置完成!${NC}"
echo ""
echo -e "${YELLOW}📋 下一步操作:${NC}"
echo "  1. 更新docker-compose.yml中的Nginx volumes，添加:"
echo "     - ./docker/nginx/certbot:/etc/letsencrypt:ro"
echo "  2. 重启Nginx容器: docker-compose restart nginx"
echo "  3. 设置自动续期（可选）: 添加cron任务执行 ./scripts/setup-ssl.sh renew"
echo ""

