COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +%F)
BUILD=$(shell echo "${BUILDNUMBER}")
NAME="windiskmgmt"
GO_LDFLAGS=-ldflags "-X main.Version=build="$(BUILD)"|commit="$(COMMIT)"|date="$(DATE)""

all: clean fmt lint vet megacheck cover

.PHONY: windows
windows:
	CGO_ENABLED=1 CC=/data/junk/bin/x86_64-w64-mingw32-gcc GOOS=windows go build -o windiskmgmt.exe cmd/cmd.go

.PHONY: linux
linux:
	go build -o windiskmgmt cmd/cmd.go

.PHONY: clean
clean:
	rm -v ${NAME} | tee /dev/stderr ; rm -v ${NAME}.exe | tee /dev/stderr ; rm -rfv data | tee /dev/stderr ; rm -v diskmgmt.log | tee /dev/stderr

.PHONY: fmt
fmt:
	gofmt -s -l . | grep -v '.pb.go' | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint:
	golint ./... | grep -v '.pb.go' | grep -v vendor | tee /dev/stderr

.PHONY: vet
vet:
	go vet $(shell go list ./... | grep -v vendor) | grep -v '.pb.go' | tee /dev/stderr

.PHONY: megacheck
megacheck:
	megacheck $(shell go list ./... | grep -v vendor) | grep -v '.pb.go' | tee /dev/stderr

.PHONY: cover
cover: ## Runs go test with coverage
	@echo "" > coverage.txt
	@for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=atomic "$$d"; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done;

