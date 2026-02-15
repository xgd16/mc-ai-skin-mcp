.PHONY: build build-dev run clean test

BINARY := mc-ai-skin
# -s -w: strip 符号表与调试信息，减小体积
# -trimpath: 移除路径信息，便于分发
BUILD_FLAGS := -ldflags "-s -w" -trimpath

build:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BINARY) .

build-dev:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test ./...
