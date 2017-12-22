FROM golang:alpine

ENV FFMPEG_VERSION=3.4.1

WORKDIR /tmp/ffmpeg

RUN apk add --update coreutils build-base git curl nasm tar bzip2 \
  zlib-dev libressl-dev yasm-dev lame-dev libogg-dev x264-dev libvpx-dev libvorbis-dev x265-dev freetype-dev libass-dev libwebp-dev rtmpdump-dev libtheora-dev opus-dev && \
  DIR=$(mktemp -d) && cd ${DIR} && \
  curl -s http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz | tar zxvf - -C . && \
  cd ffmpeg-${FFMPEG_VERSION} && \
  ./configure \
  --enable-version3 --enable-gpl --enable-nonfree --enable-small --enable-libmp3lame --enable-libx264 --enable-libx265 --enable-libvpx --enable-libtheora --enable-libvorbis --enable-libopus --enable-libass --enable-libwebp --enable-librtmp --enable-postproc --enable-avresample --enable-libfreetype --disable-debug && \
  make -j4 && \
  make install && \
  make distclean && \
  rm -rf ${DIR} && \
  apk del coreutils build-base curl tar bzip2 x264 libressl nasm && rm -rf /var/cache/apk/*

### Start of main app ###
WORKDIR /go/src/app

COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
