#!/usr/bin/env bash

set -e
set -x

# Dynamically find kernel headers. Adjust this line according to your needs.
KERNEL_HEADERS=$(find /usr/src -name "linux-headers-$(uname -r)" -type d | head -n 1)

# Include directory for bpf_helpers.h. This might need adjustments.
BPF_HELPERS_DIR="${KERNEL_HEADERS}/tools/bpf/resolve_btfids/libbpf/include/"

# Run bpf2go with dynamic include paths
go run github.com/cilium/ebpf/cmd/bpf2go -target amd64 bpf \
    ../../bpf/uretprobe.c -- \
    -I"${KERNEL_HEADERS}" \
    -I"${BPF_HELPERS_DIR}" \
    -I../../bpf/headers
