help: doc

doc:
	./scripts/gendoc.sh

gen:
	go generate ./...

test: gen
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
