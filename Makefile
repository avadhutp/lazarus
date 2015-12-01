PKGS := \
	ui \
	geddit \
	. \

SOURCES := $(foreach pkg, $(PKGS), $(wildcard $(pkg)/*.go))

lint: $(SOURCES)
	@echo Linting lazarus sources...
	@go get -u github.com/golang/lint/golint
	@$(foreach src, $(SOURCES), golint $(src);)