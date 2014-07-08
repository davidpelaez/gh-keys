fmt:
	@gofmt -w gh-keys

build: fmt
	@cd gh-keys && go build -o ../binaries/gh-keys

#-t $$(cat token.private)
run: build
	@binaries/gh-keys -i -v

test:
	@mkdir -p test
	@rm -rf test/keys test/config/keys test/main.test || :
	@cd test && go test ../gh-keys/*.go -c && ./main.test

# 
coverage:
	@cd test && go test -cover -c ../gh-keys/*.go && ./main.test 2>&1 | grep -o -E '(coverage).*' | cat

.PHONY: fmt build run test coverage