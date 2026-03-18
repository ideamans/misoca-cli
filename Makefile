BINARY_NAME=misoca

build:
	go build -o $(BINARY_NAME) ./cmd/main.go

install: build
	cp $(BINARY_NAME) /usr/local/bin/

clean:
	rm -f $(BINARY_NAME)

.PHONY: build install clean
