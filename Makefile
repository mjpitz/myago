help: doc

doc:
	./scripts/gendoc.sh

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
