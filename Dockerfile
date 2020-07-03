FROM golang:1.13.12-alpine3.12

RUN apk -U add wget bc build-base gawk xorriso elfutils-dev openssl openssl-dev bison flex ncurses-dev xz autoconf automake docbook2x alpine-sdk
RUN apk -U add linux-headers perl
RUN apk -U add rsync git
RUN apk -U add argp-standalone
RUN apk -U add xz-dev libmnl-dev libnftnl-dev libnfnetlink-dev

COPY . /build
RUN mv /build/build.sh /build.sh
RUN chmod +x /build.sh

ENTRYPOINT ["/build.sh"]
