OUT:=$(PWD)/bin
PATH:=$(PATH):$(OUT)

TOOLS:=$(OUT)

install_tools:
	GOBIN=$(OUT) go install github.com/vektra/mockery/v3@v3.0.2
.PHONY: install_tools


gen_mocks:
	$(TOOLS)/mockery
.PHONY: gen_mocks

gomod:
	go mod tidy
.PHONY: gomod