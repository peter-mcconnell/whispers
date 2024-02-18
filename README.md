whispers
========

Dependencies:

- golang
- llvm-strip
- clang
- linux headers

Build:

```sh
make docker-build
```

Usage:

```sh
Usage of ./bin/whispers:
  -binPath string
    	Path to the binary (default "/lib/x86_64-linux-gnu/libpam.so.0")
  -symbol string
    	Symbol to target (default "pam_get_authtok")
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
ldd <directory given from previous command> | grep pam


# show auth symbols of libpam
readelf -Ws /lib/x86_64-linux-gnu/libpam.so.0 | grep auth

# manually list probe events
bpftrace -e 'uprobe:/lib/x86_64-linux-gnu/libpam.so.0:pam_get_authtok { printf("pam_get_authtok called\n"); }'
```

## Credits

 - I grabbed the libpam definitions from https://github.com/citronneur/pamspy
