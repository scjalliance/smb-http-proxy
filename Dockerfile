FROM golang:latest
MAINTAINER dusty.wilson@scjalliance.com

RUN apt-get update \
    && apt-get install -y \
       cifs-utils \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 80
VOLUME /mnt/smb

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
RUN chmod 0755 /go/src/app/run.sh
RUN go-wrapper download
RUN go-wrapper install

CMD ["/go/src/app/run.sh"]
