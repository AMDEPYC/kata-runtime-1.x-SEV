#
# Copyright (c) 2018 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

From docker.io/centos:7



RUN yum -y update && yum install -y \
    autoconf \
    automake \
    binutils \
    chrony \
    coreutils \
    curl \
    gcc \
    gcc-c++ \
    git \
    glibc-common \
    glibc-devel \
    glibc-headers \
    glibc-static \
    glibc-utils \
    libseccomp \
    libseccomp-devel \
    libstdc++-devel \
    libstdc++-static \
    m4 \
    make \
    sed \
    tar \
    vim \
    which

# install cmake because centos7's cmake is too old

RUN pushd /root;     curl -sLO https://github.com/Kitware/CMake/releases/download/v3.15.3/cmake-3.15.3.tar.gz; 	tar -zxf cmake-3.15.3.tar.gz; 	cd cmake-3.15.3; 	./bootstrap > /dev/null 2>&1; 	make > /dev/null 2>&1; 	make install > /dev/null 2>&1; 	popd


RUN pushd /root;     curl -sLO https://www.musl-libc.org/releases/musl-1.1.23.tar.gz; tar -zxf musl-1.1.23.tar.gz; 	cd musl-1.1.23; 	sed -i "s/^ARCH = .*/ARCH = x86_64/g" dist/config.mak; 	./configure > /dev/null 2>&1; 	make > /dev/null 2>&1; 	make install > /dev/null 2>&1; 	echo "/usr/local/musl/lib" > /etc/ld-musl-x86_64.path; 	popd
ENV PATH=$PATH:/usr/local/musl/bin

# This will install the proper golang to build Kata components

RUN cd /tmp ; curl -OL https://storage.googleapis.com/golang/go1.11.10.linux-amd64.tar.gz
RUN tar -C /usr/ -xzf /tmp/go1.11.10.linux-amd64.tar.gz
ENV GOROOT=/usr/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin


RUN curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSLf --output /tmp/rust-init;     chmod a+x /tmp/rust-init; 	export http_proxy=; 	export https_proxy=; 	/tmp/rust-init -y
RUN . /root/.cargo/env;     export http_proxy=; 	export https_proxy=; 	cargo install cargo-when; 	rustup toolchain install 1.38.0; 	rustup default 1.38.0; 	rustup target install x86_64-unknown-linux-musl
RUN ln -sf /usr/bin/g++ /bin/musl-g++

