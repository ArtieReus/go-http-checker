PKG_NAME:=github.com/ArtieReus/go-http-checker
BUILD_DIR:=bin
BINARY:=$(BUILD_DIR)/go-http-checker
BINARY-WIN:=$(BINARY)-win
LDFLAGS:=-s -w -X github.com/ArtieReus/go-http-checker/version.GITCOMMIT=`git rev-parse --short HEAD`


.PHONY: help
help:
	@echo
	@echo "Available targets:"
	@echo "  * build             - build the binary, output to $(BINARY)"
	@echo "  * build-wi          - build the binary for windows, output to $(BINARY)"
	@echo "  * metalint          - run metalint checks"


.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) -ldflags="$(LDFLAGS)" $(PKG_NAME)
	GOOS=windows GOARCH=amd64 go build -o $(BINARY-WIN) -ldflags="$(LDFLAGS)" $(PKG_NAME)

.PHONY: metalint
metalint:
	gometalinter --vendor --disable-all -E goimports -E staticcheck -E ineffassign -E gosec --deadline=60s ./...