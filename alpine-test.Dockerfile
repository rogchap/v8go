ARG GOLANG_IMAGE=golang:1.19.3-alpine3.16
FROM $GOLANG_IMAGE

RUN apk add --update \
  alpine-sdk \
  g++ \
  gcc \
  git \
  libc-dev \
  glib-dev \
  libstdc++

ENV CGO_CXXFLAGS="-Werror"
WORKDIR v8go
COPY . ./

CMD go test --tags muslgcc -v -coverprofile c.out ./... && go tool cover -html=c.out -o /dev/stdout
