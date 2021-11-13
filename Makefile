help: doc

doc:
	./scripts/gendoc.sh

lint:
	golangci-lint run --max-issues-per-linter 0 --max-same-issues 0

gen:
	go generate ./...

test: gen
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
