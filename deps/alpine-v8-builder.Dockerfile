#
# Building V8 for alpine
#
FROM alpine:3.16.3 as v8

COPY depot_tools ./depot_tools
COPY include ./include
COPY v8 ./v8
COPY .gclient compile_v8.py ./
COPY alpine_x86_64 ./alpine_x86_64

RUN \
  apk add --update --virtual .v8-build-dependencies \
    bash \
    curl \
    g++ \
    gcc \
    glib-dev \
    icu-dev \
    libstdc++ \
    linux-headers \
    make \
    ninja \
    python3 \
    tar \
    xz \
  && cp alpine_x86_64/gn depot_tools/gn \
  && ln -s /usr/bin/python3 /usr/bin/python \
  # Compile V8
  && ./compile_v8.py --no-clang --arch x86_64
