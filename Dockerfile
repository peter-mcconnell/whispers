FROM ubuntu:22.04 AS base
ARG KERNEL
RUN apt-get update -yq && \
		apt-get install -yq curl git make clang llvm linux-headers-$(KERNEL) libbpf-dev
ARG GOARCH
ENV GOARCH=$GOARCH
RUN curl -O -L https://go.dev/dl/go1.22.2.linux-${GOARCH}.tar.gz && \
		tar -C /usr/local -xzf go1.22.2.linux-${GOARCH}.tar.gz && \
		rm -rf /var/lib/apt/lists/*

WORKDIR /src

ENV PATH=$PATH:/usr/local/go/bin

# dependencies
COPY go.mod /src
COPY go.sum /src
RUN go mod tidy

FROM ubuntu:22.04 AS build
COPY --from=base / /
COPY . /src
WORKDIR /src
ENV PATH=$PATH:/usr/local/go/bin
ARG GOARCH
ENV GOARCH=$GOARCH
RUN make whispers GOARCH=$GOARCH

ARG KERNEL
FROM ubuntu:22.04 AS sshserver
ENV DEBIAN_FRONTEND=noninteractive
ENV PATH=$PATH:/usr/local/go/bin
RUN apt-get update -yq && \
		apt-get install -yq openssh-server && \
		# some debug utilities, to aide exploration
		apt-get install -yq binutils bpftrace systemtap systemtap-sdt-dev linux-headers-$(KERNEL) vim && \
		mkdir -p /var/run/sshd && \
		echo 'root:pass' | chpasswd && \
		sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
		echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config
COPY --from=build /src/bin/whispers /bin/whispers
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
