GO_CMD=go
GO_TEST=${GO_CMD} test
GO_BUILD=${GO_CMD} build
RELEASER_CMD=goreleaser
RELEASE=${RELEASER_CMD} release
LINTER_CMD=golangci-lint
LINT=${LINTER_CMD} run
BINARY_NAME=go-lol
TARGET=dist
BIN=${TARGET}/bin

.PHONY: run all snapshot test lint clean

run:
	${GO_CMD} run cmd/go-lol/main.go

all: lint test build

snapshot: lint test snapshot
	${RELEASE} --snapshot --rm-dist

build: ${BIN}/${BINARY_NAME}

${BIN}/${BINARY_NAME}:
	mkdir -p ${BIN}
	${GO_BUILD} -o ${BIN}/${BINARY_NAME} cmd/go-lol/main.go

test:
	${GO_TEST} -race ./...

lint:
	${LINT} ./...

clean:
	go clean
	rm -rf -- ${TARGET}
	rm -f -- coverage.out

requirements:
	go install github.com/goreleaser/goreleaser@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
