fmt:
	@gofmt -w gh-keys

build: fmt
	@cd gh-keys && go build -o ../binaries/gh-keys

#-t $$(cat token.private)
run: build
	@binaries/gh-keys -i -v

test:
	@mkdir -p test
	@rm -rf test/*
	@cd test && go test ../gh-keys/*.go -c && ./main.test

coverage:
	@cd gh-keys && go test -cover 2>&1 | grep -o -E '(coverage).*' | cat

.PHONY: fmt build run test coverage