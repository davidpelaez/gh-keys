bindata:
	@go-bindata -prefix data data/...

build: bindata
	@go build

run: bindata
	@go run *.go davidpelaez

fmt:
	echo pending...