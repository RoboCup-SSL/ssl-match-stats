.PHONY: all test install proto local-match-stats-db remote-flyway

all: install

test:
	go test ./...

install:
	go install -v ./...

proto:
	tools/generateProto.sh

local-match-stats-db:
	ssl-match-stats-db -parallel=16 -sqlDbSource="postgres://ssl_match_stats:ssl_match_stats@localhost:5432/ssl_match_stats?sslmode=disable" -tournament=Test -division=DivA stats/*.bin

remote-flyway:
	docker run --net host -v $(shell pwd)/sql:/flyway/sql -v $(shell pwd)/flyway.conf:/flyway/flyway.conf flyway/flyway:6.0.8 migrate
