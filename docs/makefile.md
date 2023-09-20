## Makefile 的简单应用

```makefile
# 用于跨平台Go编译的Makefile

# Go编译器设置
GO := go
GOARCH_arm64 := arm64
GOARCH_amd64 := amd64

# 输出目录
OUTPUT_DIR := build
OUTPUT_NAME := $(CURDIR)
LINUX_ARM64_OUTPUT := $(OUTPUT_DIR)/linux_arm64
LINUX_AMD64_OUTPUT := $(OUTPUT_DIR)/linux_amd64
WINDOWS_AMD64_OUTPUT := $(OUTPUT_DIR)/windows_amd64

# Go源文件
SOURCE := main.go

.PHONY: all linux_arm64 linux_amd64 windows_amd64 clean

all: linux_arm64 linux_amd64 windows_amd64

linux_arm64:
	@echo "Compiling for Linux arm64..."
	@mkdir -p $(LINUX_ARM64_OUTPUT)
	GOOS=linux GOARCH=$(GOARCH_arm64) $(GO) build -o $(LINUX_ARM64_OUTPUT)/linux_arm64 -tags "netgo" -ldflags "-w -s" $(SOURCE)
	@echo "Done"

linux_amd64:
	@echo "Compiling for Linux amd64..."
	@mkdir -p $(LINUX_AMD64_OUTPUT)
	GOOS=linux GOARCH=$(GOARCH_amd64) $(GO) build -o $(LINUX_AMD64_OUTPUT)/linux_amd64 -tags "netgo" -ldflags "-w -s" $(SOURCE)
	@echo "Done"

windows_amd64:
	@echo "Compiling for Windows amd64..."
	@mkdir -p $(WINDOWS_AMD64_OUTPUT)
	GOOS=windows GOARCH=$(GOARCH_amd64) $(GO) build -o $(WINDOWS_AMD64_OUTPUT)/main.exe -tags "netgo" -ldflags "-w -s" $(SOURCE)
	@echo "Done"

clean:
	@echo "Cleaning..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Done"
```
