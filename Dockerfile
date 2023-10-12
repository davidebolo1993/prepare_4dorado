FROM ubuntu:latest
LABEL description="prepare_4dorado"
LABEL base_image="ubuntu:latest"
LABEL software="split_by_channel"
LABEL about.home="https://github.com/davidebolo1993/prepare_4dorado"
LABEL about.license="GPLv3"

ARG DEBIAN_FRONTEND=noninteractive
#install basic libraries and golang

WORKDIR /opt

RUN apt-get update

RUN apt-get -y install build-essential \
	software-properties-common \
	wget curl git \
	bzip2 libbz2-dev \
	zlib1g zlib1g-dev \
	liblzma-dev \
	libssl-dev \
	libncurses5-dev \
	libz-dev \
	python3-distutils python3-dev python3-pip \ 
	libjemalloc-dev \
	cmake make g++ \
	libhts-dev \
	libzstd-dev \
	autoconf \
	libatomic-ops-dev \
	pkg-config \
	pigz \
	clang-14 \ 
	libomp5 libomp-dev libssl-dev libssl3 pkg-config \
	zip unzip

#install golang
RUN add-apt-repository ppa:longsleep/golang-backports

RUN git clone https://github.com/davidebolo1993/prepare_4dorado \
	&& cd prepare_4dorado \
	&& go mod init split_by_channel \
	&& go mod tidy \
	&& go build split_by_channel

ENV PATH /opt/prepare_4dorado:$PATH
