.PHONY: swagger clean build-local build-linux build-mac run dev test install-tools

# 安装开发工具
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest

# 生成 Swagger 文档
swagger:
	swag init -g swagger/swagger.go -d web,../pkg,./iotcore -o web/docs --parseInternal

# 清理生成的文件
clean:
	rm -f docs/docs.go docs/swagger.json docs/swagger.yaml
	rm -rf target

# 本地 windows 平台构建
build-local: swagger
	go build -o target/jh-iot.exe cmd/main.go

# linux 平台构建
build-linux: swagger
	GOOS=linux GOARCH=amd64 go build -o target/jh-iot cmd/main.go

# macOS 平台构建
build-mac: swagger
	GOOS=darwin GOARCH=amd64 go build -o target/jh-iot cmd/main.go

# 运行应用（开发模式）
dev: swagger
	go run cmd/main.go

# 运行测试
test:
	go test ./...

# 格式化代码
fmt:
	goimports-reviser -project-name github.com/zhuofanxu/axb -format ./...

# 代码检查
vet:
	go vet ./...

# 下载依赖
deps:
	go mod download
	go mod tidy

# 完整的开发环境设置
setup: install-tools deps swagger

# 帮助信息
help:
	@echo Available make commands:
	@echo   swagger      - Generate Swagger docs
	@echo   build-local  - Build app for Windows (target/jh-iot.exe)
	@echo   build-linux  - Build app for Linux (target/jh-iot)
	@echo   build-mac    - Build app for macOS (target/jh-iot)
	@echo   dev          - Run app in dev mode
	@echo   test         - Run tests
	@echo   clean        - Clean generated files
	@echo   fmt          - Format code
	@echo   vet          - Code vetting
	@echo   deps         - Download dependencies
	@echo   setup        - Init dev environment
	@echo   install-tools - Install dev tools