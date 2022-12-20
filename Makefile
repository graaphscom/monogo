.PHONY: start-db lint-all test-all

modules := asciiui compoas crawler dbmigrat fa2ts rspns

start-db:
	docker run --rm -e POSTGRES_PASSWORD=dbmigrat -e POSTGRES_USER=dbmigrat -d -p 5432:5432 postgres:15.1

lint-all:
	golangci-lint run $(foreach module,$(modules),./$(module)/...)

test-all:
	DBMIGRAT_TEST_DB_URL="postgres://dbmigrat:dbmigrat@localhost:5432/dbmigrat?sslmode=disable" \
	go test github.com/graaphscom/monogo/...
