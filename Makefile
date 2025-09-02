.PHONY: build run run-bin clean

build:
	mkdir -p bin
	# build only the current package (single main)
	go build -o bin/go-streamer .

run: 
	go run .

run-bin: build
	./bin/go-streamer

clean:
	go clean
	rm -f bin/go-streamer
