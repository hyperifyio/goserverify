.PHONY: build clean tidy

GOSERVERIFY_SOURCES := \
    ./cmd/goserverify/main.go

all: build

build: goserverify

tidy:
	go mod tidy

goserverify: $(GOSERVERIFY_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o goserverify ./cmd/goserverify
	chmod 700 ./goserverify

test: cert.pem client-cert.pem ca.pem
	go test -v ./...

clean:
	rm -f goserverify
