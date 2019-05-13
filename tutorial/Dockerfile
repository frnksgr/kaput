FROM ubuntu:18.04

LABEL maintainer="Frank Seeger <frnksgr@gmail.com>"

# Install all OS dependencies
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get upgrade -yq \
    &&  apt-get install -yq --no-install-recommends \
        python3 \
        python3-dev \
        python3-pip \
        python3-setuptools \
        coreutils \
        netcat-openbsd \
        iproute2 \
    &&  rm -rf /var/lib/apt/lists/*

RUN pip3 install \
        jupyter \
    &&  mkdir tutorial

WORKDIR /tutorial
CMD jupyter notebook --no-browser --allow-root --ip $(hostname -I)