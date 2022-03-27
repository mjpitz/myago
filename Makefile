help: doc

doc:
	./scripts/gendoc.sh

lint:
	./scripts/nogoogle.py
	golangci-lint run --max-issues-per-linter 0 --max-same-issues 0

gen:
	go generate ./...

test: gen
	#go test -v -race -coverprofile=.coverprofile -covermode=atomic ./...
	go test -v -coverprofile=.coverprofile -covermode=atomic ./...

legal: .legal
.legal:
	addlicense -f ./templates/legal/header.txt -skip yaml -skip yml .
