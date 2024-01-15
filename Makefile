all:
	CGO_ENABLED=1 go build -o totp

install: all
	sudo cp totp /usr/local/bin/totp

clean:
	rm totp
