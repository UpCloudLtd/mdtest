GO       = go
MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)

ifeq ($(shell $(GO) env GOOS),windows)
    EXT = ".exe"
else
    EXT = ""
endif

BIN_DIR              = $(CURDIR)/bin
CLI_BIN              = mdtest$(EXT)

.PHONY: build
build: $(BIN_DIR)
	$(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/globals.Version=$(VERSION) -X $(MODULE)/globals.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(CLI_BIN) main.go

$(BIN_DIR):
	@mkdir -p $@

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)

.PHONY: version
version:
	@echo $(VERSION)
