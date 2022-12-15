.PHONY: start-db check-fmt test test-ci

modules := asciiui compoas crawler dbmigrat fa2ts rspns

start-db:
	docker run -e POSTGRES_PASSWORD=dbmigrat -e POSTGRES_USER=dbmigrat -d -p 5432:5432 postgres:13.3

check-fmt:
	DIFF=$$(gofmt -d .);echo "$${DIFF}";test -z "$${DIFF}"

test-ci:
	go test -coverprofile=coverage.out -covermode=set

lint-all:
	for module in $(modules); do golangci-lint run $$module; done

test-all:
	go test github.com/graaphscom/monogo/...
	#for module in $(modules); do go test ./$$module -covermode=set -failfast; done
