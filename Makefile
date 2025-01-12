.PHONY: all test install proto local-match-stats-db remote-flyway

all: install

test:
	go test ./...

install:
	go install -v ./...

proto:
	tools/generateProto.sh

flyway-migrate:
	docker run --net host -v $(shell pwd)/sql/flyway:/flyway/sql -v $(shell pwd)/flyway.conf:/flyway/flyway.conf flyway/flyway:10 migrate

flyway-repair:
	docker run --net host -v $(shell pwd)/sql/flyway:/flyway/sql -v $(shell pwd)/flyway.conf:/flyway/flyway.conf flyway/flyway:10 repair

update-go:
	go get -v -u all

update: update-go proto