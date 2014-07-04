build:
	@gofmt -w gh-keys
	@cd gh-keys && go build -o ../binaries/gh-keys

#-t $$(cat token.private)
run: bindata
	@go run gh-keys/*.go -i -v

test:
	@go test gh-keys/*.go