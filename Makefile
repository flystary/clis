#make

.PHONY: build
build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-w -s"  -o  opskey main.go
