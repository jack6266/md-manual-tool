# Makefile for md-manual-tool

.PHONY: build clean test help

# 默认目标
all: build

# 构建程序
build:
	@echo "Building md-manual-tool..."
	go build -o md-manual-tool.exe
	@echo "Build completed successfully!"

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	@if exist md-manual-tool.exe del md-manual-tool.exe
	@echo "Clean completed!"

# 运行测试
test:
	@echo "Running tests..."
	go test ./pkg/...

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./pkg/...

# 检查代码
vet:
	@echo "Checking code..."
	go vet ./pkg/...

# 构建并运行
run: build
	@echo "Running md-manual-tool..."
	./md-manual-tool.exe

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build  - Build the application"
	@echo "  clean  - Clean build files"
	@echo "  test   - Run tests"
	@echo "  fmt    - Format code"
	@echo "  vet    - Check code"
	@echo "  run    - Build and run"
	@echo "  help   - Show this help" 