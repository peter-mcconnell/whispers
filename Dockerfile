FROM ubuntu:22.04 AS build
RUN apt-get update -yq && \
		apt-get install -yq curl clang llvm && \
		apt-get install -yq bpftrace vim systemtap systemtap-sdt-dev linux-headers-$(uname -r) && \
		curl -O -L https://go.dev/dl/go1.21.7.linux-amd64.tar.gz && \
		tar -C /usr/local -xzf go1.21.7.linux-amd64.tar.gz && \
		rm -rf /var/lib/apt/lists/*

WORKDIR /src

ENV PATH=$PATH:/usr/local/go/bin

# dependencies
COPY go.mod /src
COPY go.sum /src
RUN go mod tidy

# build
COPY . /src
RUN make whispers

FROM ubuntu:22.04 AS sshserver
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -yq && \
		apt-get install -yq openssh-server && \
		apt-get install -yq binutils bpftrace systemtap systemtap-sdt-dev linux-headers-$(uname -r) vim && \
		mkdir -p /var/run/sshd && \
		echo 'root:pass' | chpasswd && \
		sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
		echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config
COPY --from=build /src/bin/whispers /bin/whispers
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]