VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

BINARY_NAME_LINUX = ac-tts-${VERSION}.out
BINARY_NAME_WIN64 = ac-tts-${VERSION}.exe

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o ${BINARY_NAME_LINUX} main.go
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags "-X main.Version=$(VERSION)" -o ${BINARY_NAME_WIN64} main.go

run:
	./${BINARY_NAME}
clean:
	go clean
	rm -f $(BINARY_NAME_LINUX) $(BINARY_NAME_WIN64)

deps:
	go get -v -t ./...