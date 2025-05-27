APP_NAME := gt7_car_sales

.PHONY: all clean

all: macos windows linux

macos:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -o bin/$(APP_NAME)_macos .

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -o bin/$(APP_NAME)_windows.exe .

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -o bin/$(APP_NAME)_linux .

clean:
	rm -rf bin/