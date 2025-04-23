.PHONY: deps build-linux build-windows clean all

deps:
	go mod download

build-linux:
	GOOS=linux GOARCH=amd64 go build -o rg-downloader main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o rg-downloader.exe main.go

clean:
	rm -f rg-downloader rg-downloader.exe

all: clean deps build-linux build-windows
