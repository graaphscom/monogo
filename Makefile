.PHONY: start-db check-fmt test test-ci

modules := asciiui compoas crawler dbmigrat fa2ts rspns

start-db:
	docker run --rm -e POSTGRES_PASSWORD=monogo -e POSTGRES_USER=monogo -d -p 5432:5432 postgres:13.3

lint-all:
	golangci-lint run $(foreach module,$(modules),./$(module)/...)

test-all:
	go test github.com/graaphscom/monogo/...
