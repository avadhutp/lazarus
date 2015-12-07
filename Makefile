PKGS := \
	ui \
	geddit \
	. \

SOURCES := $(foreach pkg, $(PKGS), $(wildcard $(pkg)/*.go))

lint: $(SOURCES)
	@echo Linting lazarus sources...
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/GeertJohan/fgt
	@$(foreach src, $(SOURCES), fgt golint ./... || exit;)

test: 
	@go test -v ./...
