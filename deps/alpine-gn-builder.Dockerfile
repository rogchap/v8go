#
# Build GN for alpine (this is a build dependency)
#
FROM alpine:3.16.3 as gn-builder

# This is the GN commit that we want to build.
# Most commits will probably build just fine but this happened to be the latest commit when I did this.
ARG GN_COMMIT=1c4151ff5c1d6fbf7fa800b8d4bb34d3abc03a41

RUN \
  apk add --update --virtual .gn-build-dependencies \
    alpine-sdk \
    binutils-gold \
    clang \
    curl \
    git \
    llvm12 \
    ninja \
    python3 \
    tar \
    xz \
  # Quick fixes: we need the LLVM tooling in $PATH, and we also have to use gold instead of ld.
  && PATH=$PATH:/usr/lib/llvm12/bin \
  && cp -f /usr/bin/ld.gold /usr/bin/ld \
  # Clone and build gn
  && git clone https://gn.googlesource.com/gn /tmp/gn \
  && git -C /tmp/gn checkout ${GN_COMMIT} \
  && cd /tmp/gn \
  && python3 build/gen.py \
  && ninja -C out \
  && cp -f /tmp/gn/out/gn /usr/local/bin/gn