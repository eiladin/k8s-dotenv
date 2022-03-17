default: help
## This help screen.
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = $$1; \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			gsub("\\\\", "", helpCommand); \
			gsub(":+$$", "", helpCommand); \
			printf "  \x1b[36;01m%-30s\x1b[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"


.PHONY: build
## run goreleaser in snapshot mode
build:
	@goreleaser release --snapshot --rm-dist  --skip-validate --skip-publish

.PHONY: test
## generate coverage profile
test:
	@go test ./... -coverprofile=coverage.out -covermode=atomic

.PHONY: cover
## open cover profile in browser
cover:
	@go tool cover -html=coverage.out

.PHONY: docs
## generate docs
docs:
	@go run main.go doc