BINARY_NAME_LINUX=ac-tts.out
BINARY_NAME_WIN64-ac-tts.exe

build:
	go build -o ${BINARY_NAME} main.go
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o ${BINARY_NAME}

run:
	./${BINARY_NAME}
clean:
	go clean
	rm ${BINARY_NAME}

deps:
	go get -v -t ./...