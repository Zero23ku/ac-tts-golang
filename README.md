# AC - TTS

How to build

```
$go mod tidy
$go build -o app
```

How to build for Windows under Linux

```
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o app.exe
```



<a href='https://ko-fi.com/I3I41O6OUD' target='_blank'><img height='36' style='border:0px;height:36px;' src='https://storage.ko-fi.com/cdn/kofi4.png?v=6' border='0' alt='Buy Me a Coffee at ko-fi.com' /></a>