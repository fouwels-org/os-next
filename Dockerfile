FROM alpine:latest

RUN apk -U add wget bc build-base gawk xorriso elfutils-dev openssl openssl-dev bison flex ncurses-dev xz autoconf automake docbook2x alpine-sdk
RUN apk -U add linux-headers perl
RUN apk -U add rsync git
RUN apk -U add argp-standalone
RUN apk -U add xz-dev libmnl-dev libnftnl-dev libnfnetlink-dev

WORKDIR /build
COPY . /build

RUN chmod +x build.sh

#ENTRYPOINT ["./build.sh"]
#CMD [""]
