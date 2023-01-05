.PHONY: http
## http :
http:
	@./scripts/api-http.sh

.PHONY: format
## format :format project code
format:
	@echo FORMAT PROJECT CODE...
	@go fmt ./...

.PHONY: gofmt
## gofmt :enforce a stricter format than gofmt, while being backwards compatible.
gofmt:
	@echo FOFUMPT PROJECT CODE...
	@go install mvdan.cc/gofumpt@latest
	@gofumpt -l -w .

.PHONY: tidy
## tidy :download need mod && clean unused mod && upgrade mod
tidy:
	@echo TIDYING CODE...
	@go mod tidy -compat=1.18

.PHONY: help
## help: prints all cmdline of help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: lint
## lint :go linters aggregator
lint:
	@echo LINTING CODE...
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	@if [ "`golangci-lint run ./... | tee /dev/stderr`" ]; then \
  		echo "^ - Lint errors!" && echo && exit 1; \
  	fi

.PHONY: compile
## compile :compile oms backend server code
compile:
	@echo COMPILE CODE...
	@go build -o ddd-johnny ./cmd/main.go ./cmd/http.go

.PHONY: serve
## serve :start oms backend server
serve:
	@echo START SERVER...
	./ddd-johnny

.PHONY: reload
## reload :only start oms backend server
reload: compile serve

.PHONY: build
## build :build oms backend server
build: tidy http gofmt format lint compile

.PHONY: release
## release :release oms backend server
release: tidy http compile

.PHONY: all
## all :default run all cmdline for server
all: tidy http gofmt format lint compile serve
