# --------
# Stage 1: Build
# -------
FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /go/src/github.com/scjalliance/smb-http-proxy
COPY . .

# Disable CGO to make sure we don't rely on libc
ENV CGO_ENABLED=0

# Exclude debugging symbols and set the netgo tag for Go-based DNS resolution
ENV BUILD_FLAGS="-v -a -ldflags '-d -s -w' -tags netgo"

RUN go-wrapper download
RUN go-wrapper install

# --------
# Stage 2: Release
# --------
FROM debian

VOLUME /mnt/smb
EXPOSE 80

#RUN apt-get update \
#    && apt-get install -y \
#       cifs-utils \
#    && apt-get clean \
#    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/smb-http-proxy /

WORKDIR /data
CMD ["/smb-http-proxy"]
