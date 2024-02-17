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
whispers -binPath="/lib/x86_64-linux-gnu/libpam.so.0" -symbol="pam_get_authtok"
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
