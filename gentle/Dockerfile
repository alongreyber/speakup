FROM ubuntu:16.04

RUN echo -e "*************\nGo get some coffee. This will take a while\n*************"
RUN mkdir gentle
WORKDIR gentle

RUN apt update && \
               apt install -y git && \
	       apt autoremove
ADD ext/ ./ext
ADD install* ./
# Needed to initialize kaldi submodule

RUN bash install.sh
ADD . /gentle

EXPOSE 8765

VOLUME /gentle/webdata

CMD python serve.py
