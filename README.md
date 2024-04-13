whispers
========

A binary that dumps credentials from libpam (as used by openssh, passwd and others). Written for a tech talk. Accompanying blog post: https://www.petermcconnell.com/posts/whispers/

[![asciicast](https://asciinema.org/a/641250.png)](https://asciinema.org/a/641250)

## Usage

```sh
Usage of ./bin/whispers:
  -binPath string
    	Path to the binary (default "/lib/x86_64-linux-gnu/libpam.so.0")
  -symbol string
    	Symbol to target (default "pam_get_authtok")
```


## Demo

Build a docker image with SSHD running and the whispers binary present

```sh
make docker-run
```

Now in a new terminal, ssh into the container we've just ran using port 2222:

```sh
ssh root@localhost -p 2222
```

Now in the previous terminal (where you ran `make docker-run`) exec into the container:

```sh
make docker-exec
```

And once inside the container, simply run `whispers`. Now repeat the ssh login in your other terminal
- whispers should dump out the credentials as you log in. Also try changing the password with `passwd`
and it should also capture this.


## Building

### Dependencies

- golang
- llvm
- clang
- libbpf-dev
- linux-tools-generic
- linux-headers
- make

### Build

Generate vmlinux.h:

```sh
make vmlinux
```

Build whispers

```sh
make whispers GOARCH=arm64  # or GOARCH=amd64
```


## Exploration

This repo was built for a talk on Golang and BPF. As such it is meant as an educational tool more than anything.
Given that, I'd encourage you to explore.

```sh
# run the docker image
make docker-run

# exec into the docker image
make docker-exec

# view linked libpam lib of sshd (/usr/sbin/sshd from inside the container)
ldd /usr/sbin/sshd | grep pam

# alteratively, exit the container and run ldd on the binary path on host:
docker inspect whispers -f '{{.GraphDriver.Data.MergedDir}}'
ldd <directory given from previous command>/usr/sbin/sshd | grep pam
# ^ if this seems interesting to you, check out https://www.petermcconnell.com/posts/docker-overlayfs/


# show auth symbols of libpam
readelf -Ws /lib/x86_64-linux-gnu/libpam.so.0 | grep auth

# manually list probe events
bpftrace -e 'uprobe:/lib/x86_64-linux-gnu/libpam.so.0:pam_get_authtok { printf("pam_get_authtok called\n"); }'
```

## Docker images

`docker pull pemcconnell/whispers-base:latest` - this is a simple docker image which contains the dependencies required to build whispers. It is used for CI

## Credits

 - I grabbed the libpam definitions from https://github.com/citronneur/pamspy
