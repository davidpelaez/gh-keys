bindata:
	@go-bindata -o gh-keys/bindata.go -prefix bindata bindata/...

build: bindata
	@gofmt -w gh-keys
	@cd gh-keys && go build -o ../binaries/gh-keys

run: bindata
	@go run gh-keys/*.go davidpelaez