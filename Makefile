APP=ai
VERSION=v1.0.0

build:
	go build -o $(APP) .

release:
	GOOS=darwin  GOARCH=arm64  go build -o dist/$(APP)-darwin-arm64 .
	GOOS=darwin  GOARCH=amd64  go build -o dist/$(APP)-darwin-amd64 .
	GOOS=linux   GOARCH=amd64  go build -o dist/$(APP)-linux-amd64 .
	GOOS=windows GOARCH=amd64  go build -o dist/$(APP)-windows-amd64.exe .

clean:
	rm -rf dist/ $(APP)

install:
	go build -o $(APP) .
	mkdir -p $(HOME)/.local/bin
	mv $(APP) $(HOME)/.local/bin/$(APP)