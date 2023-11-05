.PHONY: all test install proto

all: install

test:
	go test ./...

install:
	go install -v ./...

proto:
	tools/generateProto.sh

local-match-stats-db:
	ssl-match-stats-db -sqlDbSource="postgres://ssl_match_stats:ssl_match_stats@localhost:5432/ssl_match_stats?sslmode=disable" -tournament=Test -division=DivA
