GO := go
GOFLAGS := -ldflags="-s -w" -trimpath
LDFLAGS := -buildvcs=false
CGO_ENABLED := 0

.PHONY: all clean build install

all: build-windows build-linux

install:
	$(GO) mod tidy
	$(GO) mod download

build:
	$(GO) build $(GOFLAGS) -o bin/ai.exe cmd/ai/main.go

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o bin\ai-windows-amd64.exe cmd/ai/main.go

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o bin\ai-linux-amd64 cmd/ai/main.go

clean:
	rm -rf bin/