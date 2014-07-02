bindata:
	@go-bindata -o gh-keys/bindata.go -prefix bindata bindata/...

build: bindata
	@gofmt -w gh-keys
	@cd gh-keys && go build -o ../binaries/gh-keys

#-t $$(cat token.private)
run: bindata
	@go run gh-keys/*.go -i -v