uname_lower:= $(shell uname | tr '[:upper:]' '[:lower:]')
arch := $(shell uname -m)

all: clean
	CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin go build -o bin/totp_darwin_amd64
	CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin go build -o bin/totp_darwin_arm64

install: all
	sudo cp bin/totp_$(uname_lower)_$(arch) /usr/local/bin/totp

clean:
	rm -rf bin/
