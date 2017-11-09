FROM ubuntu:16.04

RUN DEBIAN_FRONTEND=noninteractive && \
	apt-get update && \
	apt-get install -y zlib1g-dev automake autoconf git \
		libtool subversion libatlas3-base ffmpeg python-pip \
		python-dev wget unzip && \
	apt-get clean

COPY install* /gentle/
COPY requirements.txt /gentle/
COPY ext/ /gentle/ext/

# This is the long step so we execute it before copying the code
WORKDIR /gentle
RUN ./install.sh

COPY . /gentle
RUN pip install .

EXPOSE 80
CMD python serve.py
