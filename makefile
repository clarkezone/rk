ifeq ($(strip $(VERSION_STRING)),)
VERSION_STRING := $(shell git rev-parse --short HEAD)
endif

.PHONY: test
test:
	go test -p 4 -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: integration-actions
integration-actions:
	go test github.com/clarkezone/rk/integration -tags="common actions" --count=1 -v -timeout 15m

.PHONY: dep
dep:
	go mod tidy

.PHONY: latest
latest:
	echo ${VERSION_STRING} > bin/latest