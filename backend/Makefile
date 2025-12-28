.PHONY: swag build run clean test

# 生成Swagger文档
swag:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# 构建应用
build: swag
	@echo "Building application..."
	@go build -o blog-server ./cmd/server

# 运行应用
run: swag
	@echo "Running application..."
	@go run ./cmd/server

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf blog-server tmp/ docs/

# 测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
lint:
	@echo "Linting code..."
	@golangci-lint run

# 安装Swagger工具
install-swag:
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest

# 帮助
help:
	@echo "Available targets:"
	@echo "  swag          - Generate Swagger documentation"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  install-swag  - Install swag tool"
	@echo "  help          - Show this help message"

