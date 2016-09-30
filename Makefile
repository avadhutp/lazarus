deps:
	@echo Getting dependencies for lazarus
	@go get github.com/mattn/gom
	@gom install

test:
	@echo Testing lazarus
	@(go list ./... | grep -v -e /vendor/ | xargs -L1 gom test -cover || exit;)

lint:
	@echo Linting lazarus sources
	@(go list ./... | grep -v -e /vendor/ | xargs -L1 golint || exit;)