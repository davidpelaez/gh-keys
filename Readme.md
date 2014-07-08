Github Auth Keys Command
========================

*Version notice: This project has been built with Go 1.3 and won't work in 1.2 because of the HTTP Client Timeout option added in 1.3. If you have problems compiling please make sure you have the proper version installed (or newer)*


Contributions are more than welcome, this is my first attempt to write a non trivial program in Go, so it may not be the most idiomatic code. Feel free to open issues with corrections to make things more idiomatic/organized, to discuss new features, etc.

gh-keys has been created to simplify access to company servers in different datacenter or cloud providers. Therefore, tests are very important (you don't want to be locked out of your server, do you?), so please keep that in mind for any pull requests.

You'll need godep to build gh-keys, if you don't have it, install it with:

```
go get github.com/tools/godep
```

With godeps and Go 1.3 installed, you can download the code and verify everything's fine with:

```
go get -d github.com/davidpelaez/gh-keys/gh-keys
cd $GOPATH/src/github.com/davidpelaez/gh-keys
godep restore
make test
```

Compile the code to `binaries/gh-keys`:

```
make build
```

## TODO
* don't save bootstrap key?
* installation instructions
* security notes
* godeps instructions
* test pkg is go-gettable
* comments on file structure/division
* create user account on first login, e.g? http://linux.die.net/man/8/pam_mkhomedir
* add version with something like `go build -ldflags "-X main.Godeps '$(git rev-parse HEAD)'"`
