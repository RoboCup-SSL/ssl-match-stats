.PHONY: all test install proto

all: install

test:
	go test ./...

install:
	go install -v ./...

proto:
	tools/generateProto.sh
