Github Auth Keys Command
========================

Make your servers SSH daemon accept up to date public keys form Github's API

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

That gives you a single binary that you can then install for sshd to use.

## Config files and locations

You can pass a configuration file using the `-c <file>` command line flag, however it's more common to skip any flags and let the binary find your `config.yml` (you can also use JSON or Toml) in one of the following locations: `/etc/ghk, $HOME/.ghk, .`

For a production deployment, the most logical location is `/etc/ghk/config.yml`

Cached keys are stored in the `keys` folder relative to the configuration file. If no configuration file is provided, they're stored relative to the `gh-keys` binary.

You can see a typical config file in the `examples` directory. To see all the options (e.g. private token for higher API limits, connectivity check endpoints, etc.) that you can control, please check `gh-keys/config.go`

## Installation

You must install the binary after compilation, create the configuration folder and ensure all the right permissions are set, for instance from the root of the repo after `make build` you could (as root):

```
cp binaries/gh-keys /usr/sbin/gh-keys
chown root:root /usr/sbin/gh-keys
chmod 0700 /usr/sbin/gh-keys
mkdir -p /etc/ghk/keys
touch /etc/ghk/config.yml
chmod 0600 /etc/ghk/config.yml
chmod 0700 /etc/ghk/keys
chown root:root -R /etc/ghk/
```

That's somewhat repetitive and in some cases some commands don't actually change anything if performed as root the first time, but it's just so you see clearly how permissions should end up applied.


Change `/etc/ghk/config.yml` to set some permissions (at least) or any other config options and then check the config summary with `/sbin/gh-keys -i` to verify it looks as desired and there are no typos, etc.

Check the expected public keys are printed when you call `gh-keys`, e.g. `gh-keys root` should print the keys of the github users that can login as root as set in the config file. You can see the keys being cached in `/etc/ghk/keys/`

If you wan to be savvy, you can try to change those files as another user and you should get `permission denied` errors on write/touch, etc.

Configure sshd to use `gh-keys` as they authorized keys command. There's a sample config in the examples folder. Edit your `sshd_config` file, typically located in `/etc/ssh/sshd_config` (notice it's sshd, not ssh), or copy the sample. The sample file is a slightly modified version of Fedora's 20 version without the comments or unnecessary empty lines. The important parts are the directives related to `AuthorizedKeysCommand`.

Check the configuration file is correct with `sshd -t`

If all checks worked as expected you can proceed to restart sshd, for example in Fedora this would be `systemctl restart sshd`. Then verify that you can now login to your server with the same key(s) that you have configured in your Github account.

Proceed to delete all authorized key files from your user's homes since those will continue to work if present, e.g. `rm $HOME/.ssh/authorized_keys` would delete that file for the current user. You can alternatively change the config file to not accept the keys in those files.

Now if you change your Github key, after the configured TTL (or 5 mins by default) you will be able to use the new key to login into the machine.

## Security notice

Anyone who can change the binary could potentially allow any desired public key to access any account, so it's very important that you check file permissions are set accordingly. The sample applies for anyone who can change your config file or keys folder. Basically you need to protect this with the same care that you would use in other sensible files like `/etc/ipsec.conf` or `/etc/sudoers`

The bootstrap key is a special (typically public) that you use to access the machine when no permissions have ben set in the configuration file. This allows for example Vagrant machines or generic AWS AMIs to have gh-keys pre-installed. *If you don't explicitely configure at least one permission in the config file the bootstrap key can be used to access to _any_ account that SSH accepts logins to.* This key default to [vagrant's public key](https://github.com/mitchellh/vagrant/tree/master/keys).

## TODO

* review all instructions and simplify if possible
* make github pages explaining benefits
* get configuration from an HTTP endpoint for centralized control
* create user account on first login, e.g http://linux.die.net/man/8/pam_mkhomedir
* add version with something like `go build -ldflags "-X main.Godeps '$(git rev-parse HEAD)'"`
