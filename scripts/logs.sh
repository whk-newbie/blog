#!/bin/bash

# 查看日志

SERVICE=${1:-all}

if [ "$SERVICE" == "all" ]; then
    echo "📝 Viewing all logs (Ctrl+C to exit)..."
    docker-compose logs -f
else
    echo "📝 Viewing logs for $SERVICE (Ctrl+C to exit)..."
    docker-compose logs -f "$SERVICE"
fi

