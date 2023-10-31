.PHONY: start-db start-redis lint-all test-all

modules := asciiui compoas crawler dbmigrat icommon rspns

start-db:
	docker run --rm -e POSTGRES_PASSWORD=dbmigrat -e POSTGRES_USER=dbmigrat -d -p 5432:5432 postgres:15.1

start-redis-stack:
	docker run --rm -p 6379:6379 -p 8001:8001 redis/redis-stack

lint-all:
	golangci-lint run $(foreach module,$(modules),./$(module)/...)

test-all:
	DBMIGRAT_TEST_DB_URL="postgres://dbmigrat:dbmigrat@localhost:5432/dbmigrat?sslmode=disable" \
	go test github.com/graaphscom/monogo/...
