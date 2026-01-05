VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

BINARY_NAME_LINUX = ac-tts-${VERSION}.out
BINARY_NAME_WIN64 = ac-tts-${VERSION}.exe
BINARY_NAME_DARWIN_AMD64 = ac-tts-${VERSION}-darwin-amd64
BINARY_NAME_DARWIN_ARM64 = ac-tts-${VERSION}-darwin-arm64

build-linux:
	CGO_ENABLED=1 go build -ldflags "-X main.Version=$(VERSION)" -o ${BINARY_NAME_LINUX} main.go

build-windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -v -x -ldflags "-X main.Version=$(VERSION)" -o ${BINARY_NAME_WIN64} main.go

build-darwin-amd64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ${BINARY_NAME_DARWIN_AMD64} main.go

build-darwin-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o ${BINARY_NAME_DARWIN_ARM64} main.go

build-darwin: build-darwin-amd64 build-darwin-arm64

build: build-linux build-windows

build-all: build-linux build-windows build-darwin

run:
	./${BINARY_NAME_LINUX}

clean:
	go clean
	rm -f $(BINARY_NAME_LINUX) $(BINARY_NAME_WIN64) $(BINARY_NAME_DARWIN_AMD64) $(BINARY_NAME_DARWIN_ARM64)

deps:
	go get -v -t ./...
