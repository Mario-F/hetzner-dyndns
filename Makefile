SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

# Run all the tests
test:
	LC_ALL=C go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=5m
.PHONY: test

# Run all the linters
lint:
	golangci-lint run ./...
	misspell -error **/*
.PHONY: lint

# Run gorelease dry
godry:
	goreleaser --snapshot --skip-publish --rm-dist
.PHONY: godry
