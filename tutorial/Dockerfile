FROM ubuntu:18.04

LABEL maintainer="Frank Seeger <frnksgr@gmail.com>"

# Install all OS dependencies
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get upgrade -yq \
    &&  apt-get install -yq --no-install-recommends \
        netcat-openbsd \
        curl \
        python3 \
        wireshark \
        iproute2 \
        tmux \
    &&  rm -rf /var/lib/apt/lists/*

RUN mkdir tutorial

WORKDIR /tutorial
