.PHONY: all test install proto local-match-stats-db remote-flyway

all: install

test:
	go test ./...

install:
	go install -v ./...

proto:
	tools/generateProto.sh

flyway-remote:
	docker run --net host -v $(shell pwd)/sql:/flyway/sql -v $(shell pwd)/flyway.conf:/flyway/flyway.conf flyway/flyway:10 migrate

flyway-local:
	docker compose up flyway