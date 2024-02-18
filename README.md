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

Useful commands:

```sh
# view linked libpam lib of sshd 
ldd /usr/sbin/sshd | grep pam

# show auth symbols of libpam
readelf -Ws /lib/x86_64-linux-gnu/libpam.so.0 | grep auth

# manually list probe events
bpftrace -e 'uprobe:/lib/x86_64-linux-gnu/libpam.so.0:pam_get_authtok { printf("pam_get_authtok called\n"); }'
```

## Credits

 - I grabbed the libpam definitions from https://github.com/citronneur/pamspy
