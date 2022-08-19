APP=vanityrender

BUILD_DIR ?= out
VENDOR_DIR = vendor

GOLANGCI_LINT_VERSION ?= v1.48.0

GO ?= go
GOLANGCI_LINT ?= golangci-lint-$(GOLANGCI_LINT_VERSION)

ifneq "$(wildcard ./vendor )" ""
    modVendor =  -mod=vendor
    ifeq (,$(findstring -mod,$(GOFLAGS)))
        export GOFLAGS := ${GOFLAGS} ${modVendor}
    endif
endif

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

.PHONY: lint
lint: bin/$(GOLANGCI_LINT) $(VENDOR_DIR)
	@bin/$(GOLANGCI_LINT) run -c .golangci.yaml

.PHONY: build
build: $(VENDOR_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(APP) cmd/*

.PHONY: test
test: test-unit

## Run unit tests
.PHONY: test-unit
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

#.PHONY: test-integration
#test-integration:
#	@echo ">> integration test"
#	@$(GO) test ./features/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./... -race --godog

bin/$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint bin/$(GOLANGCI_LINT)
