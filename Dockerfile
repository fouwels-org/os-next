FROM golang:1.15-alpine3.13

RUN apk --no-cache add wget bc build-base gawk xorriso elfutils-dev openssl openssl-dev bison flex ncurses-dev xz autoconf automake docbook2x alpine-sdk libtool asciidoc readline-dev gmp-dev
RUN apk --no-cache add linux-headers perl
RUN apk --no-cache add rsync git
RUN apk --no-cache add argp-standalone
RUN apk --no-cache add xz-dev libmnl-dev libnftnl-dev cmake libnfnetlink-dev gzip ccache diffutils util-linux libuuid util-linux-dev lvm2-dev popt popt-dev json-c json-c-dev libaio-dev upx gettext-dev
RUN apk --no-cache add openssl-libs-static --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main
RUN apk --no-cache add lvm2-static --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main
RUN apk --no-cache add device-mapper-static --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main
RUN apk --no-cache add cryptsetup e2fsprogs libpciaccess-dev

WORKDIR /tmp
RUN wget http://ftp.rpm.org/popt/releases/popt-1.x/popt-1.18.tar.gz && \
    tar -xvf popt-1.18.tar.gz && cd popt-1.18 && \
    ./configure --prefix=/usr && \
    make && \
    make install 

RUN wget https://s3.amazonaws.com/json-c_releases/releases/json-c-0.15.tar.gz && \
    tar -xvf json-c-0.15.tar.gz && cd json-c-0.15 && \
    cmake -DCMAKE_INSTALL_PREFIX=/usr -DCMAKE_BUILD_TYPE=Release -DBUILD_STATIC_LIBS=ON && \ 
    make && \
    make install

COPY . /build
COPY ./init /rebuild/init

RUN cp /build/scripts/build.sh /build.sh
RUN chmod +x /build.sh
RUN mv /build/scripts/kernel-test.sh /kernel-test.sh
RUN chmod +x /kernel-test.sh

WORKDIR /build/scripts
RUN wget https://github.com/moby/moby/raw/master/contrib/check-config.sh
RUN chmod +x check-config.sh

#ENTRYPOINT ["/build.sh"]
