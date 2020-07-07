FROM golang:1.13.12-alpine3.12

RUN apk -U add wget bc build-base gawk xorriso elfutils-dev openssl openssl-dev bison flex ncurses-dev xz autoconf automake docbook2x alpine-sdk 
RUN apk -U add linux-headers perl
RUN apk -U add rsync git
RUN apk -U add argp-standalone
RUN apk -U add xz-dev libmnl-dev libnftnl-dev libnfnetlink-dev gzip ccache

COPY . /build
RUN mv /build/build.sh /build.sh
RUN chmod +x /build.sh
RUN mv /build/kernel-test.sh /kernel-test.sh
RUN chmod +x /kernel-test.sh

ENTRYPOINT ["/build.sh"]
