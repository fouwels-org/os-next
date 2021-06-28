# SPDX-FileCopyrightText: 2020 Lagoni Engineering
# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
#
# SPDX-License-Identifier: Apache-2.0

FROM alpine:3.14.0

RUN apk --no-cache add \
    wget bc build-base gawk xorriso elfutils-dev openssl openssl-dev bison flex ncurses-dev xz autoconf automake docbook2x alpine-sdk libtool asciidoc readline-dev gmp-dev linux-headers perl rsync git argp-standalone \
    xz-dev libmnl-dev libnftnl-dev cmake libnfnetlink-dev gzip ccache diffutils util-linux libuuid util-linux-dev lvm2-dev popt popt-dev json-c json-c-dev libaio-dev upx gettext-dev openssl-libs-static lvm2-static device-mapper-static \
    cryptsetup e2fsprogs libpciaccess-dev popt go tree lz4 lzo pigz tpm2-tss tpm2-tss-dev tpm2-tss-esys tpm2-tss-fapi tpm2-tss-mu tpm2-tss-sys
RUN git config --global advice.detachedHead false bash
SHELL ["/bin/bash", "-c"]

# Dirs
ENV SRC_DIR=/build/src
ENV OUT_DIR=/build/out
RUN mkdir -p ${OUT_DIR} && mkdir -p ${SRC_DIR} && mkdir -p /rootfs
WORKDIR ${SRC_DIR}

# Kernel versions
ENV VERSION_KERNEL=5.10.1
ENV VERSION_RT=5.10.1-rt20
ENV CONFIG_KERNEL=5.10.1-rt20

# Download and patch kernel
RUN wget -q -O kernel.tar.xz https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-${VERSION_KERNEL}.tar.xz && tar -xf kernel.tar.xz
RUN wget -q -O patch-rt.xz https://cdn.kernel.org/pub/linux/kernel/projects/rt/5.10/older/patch-${VERSION_RT}.patch.xz
RUN cd linux-${VERSION_KERNEL} && xzcat ../patch-$KERNEL_RT.patch.xz | patch -p1

# Build kernel with custom config
RUN touch /initramfs.cpio
COPY config/config-${CONFIG_KERNEL} linux-${VERSION_KERNEL}/.config
RUN cd linux-${VERSION_KERNEL} && make oldconfig -j $(nproc)
RUN cd linux-${VERSION_KERNEL} && make CFLAGS="-pipe -Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $(nproc)

# Package versions
ENV VERSION_MUSL=1.2.2
ENV VERSION_DOCKER=20.10.6
ENV VERSION_BUSYBOX=1.33.1
ENV VERSION_WGTOOLS=v1.0.20210424
ENV VERSION_MICROCODE_INTEL=20210216
ENV VERSION_IPTABLES=1.8.7

# Build musl
RUN wget -q -O musl.tar.gz https://www.musl-libc.org/releases/musl-$VERSION_MUSL.tar.gz && tar -xf musl.tar.gz
RUN cd musl-${VERSION_MUSL} && ./configure --prefix=/usr && \
    make -j $(nproc)

# Build docker
RUN wget -q -O docker.tgz https://download.docker.com/linux/static/stable/x86_64/docker-$VERSION_DOCKER.tgz && tar -xf docker.tgz

# Build wg-tools
RUN git clone --depth 1 -b $VERSION_WGTOOLS https://git.zx2c4.com/wireguard-tools
RUN cd wireguard-tools/src && \
    make -j $(nproc)

# Build busybox
RUN wget -q -O busybox.tar.bz2 https://busybox.net/downloads/busybox-$VERSION_BUSYBOX.tar.bz2 && tar -xf busybox.tar.bz2
COPY config/busybox-config busybox-${VERSION_BUSYBOX}/.config
RUN cd busybox-${VERSION_BUSYBOX} && make oldconfig
RUN cd busybox-${VERSION_BUSYBOX} && \
    make -j $(nproc)

# Build iptables
RUN wget -q -O iptables.tar.bz2  https://netfilter.org/projects/iptables/files/iptables-$VERSION_IPTABLES.tar.bz2 && tar -xf iptables.tar.bz2
RUN cd iptables-${VERSION_IPTABLES} && ./configure --prefix=/ --mandir=/tmp --disable-nftables && \
    make EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $(nproc)

# Set up template rootfs
COPY template_rootfs /template_rootfs
COPY scripts/build-rootfs.sh .
RUN ./build-rootfs.sh /rootfs /template_rootfs

# Install CPU microcode to rootfs
RUN wget -q -O microcode.tar.gz https://github.com/intel/Intel-Linux-Processor-Microcode-Data-Files/archive/refs/tags/microcode-${VERSION_MICROCODE_INTEL}.tar.gz && tar -xzf microcode.tar.gz
RUN mkdir -p /lib/firmware/intel-ucode && cp -r Intel-Linux-Processor-Microcode-Data-Files-microcode-${VERSION_MICROCODE_INTEL}/intel-ucode /lib/firmware/intel-ucode

# Install modules to rootfs
RUN cd linux-${VERSION_KERNEL} && mkdir -p /rootfs && make -j $(nproc) INSTALL_MOD_PATH=/rootfs modules_install

# Install packages to rootfs
RUN cd busybox-${VERSION_BUSYBOX} && make -j $(nproc) CONFIG_PREFIX=/rootfs install
RUN cd musl-${VERSION_MUSL} && make -j $(nproc) DESTDIR=/rootfs install
RUN cd iptables-${VERSION_IPTABLES} && make DESTDIR=/rootfs install
RUN cp wireguard-tools/src/wg /rootfs/usr/sbin/wg
RUN cp docker/* /rootfs/usr/bin/

# Add alpine packages
RUN apk add --no-cache coreutils
RUN cd /bin && cp -t /rootfs/bin lsblk 
RUN cd /sbin && cp -t /rootfs/bin mke2fs
RUN cd /lib && cp -t /rootfs/lib libext2fs.so.* libcom_err.so.* libblkid.so.* libuuid.so.* libe2p.so.* libsmartcols.so.* libmount.so.* libfdisk.so.*
RUN cd /usr/lib && cp -t /rootfs/usr/lib libncursesw.so.*
RUN cp /sbin/fdisk /rootfs/bin/fdisk2 # Alias GPT aware fdisk to fdisk2 to prevent clash with busybox fdisk.

# Strip modules if specified
ARG CONFIG_MODULES=ALL
COPY config/modules .
RUN find /rootfs/lib/modules | grep "\.ko$" > ${OUT_DIR}/modules.txt
RUN if [ "${CONFIG_MODULES}" != "ALL" ]; then find /rootfs/lib/modules | grep "\.ko$" | grep -v -f ${CONFIG_MODULES} | xargs rm; fi;
RUN find /rootfs/lib/modules | grep "\.ko$" > ${OUT_DIR}/modules_selected.txt

# Optimise rootfs
RUN find /rootfs -executable -type f | xargs strip || true
RUN find /rootfs | grep ".\.la$" | xargs rm || true
RUN find /rootfs | grep ".\.a$" | xargs rm || true
RUN find /rootfs | grep ".\.o$" | xargs rm || true

# Cache go/init dependencies
COPY init/go.mod init/go.mod
COPY init/go.sum init/go.sum
RUN cd init && go mod download

# Build go/init into rootfs
COPY init init
RUN cd init && go build -ldflags "-s -w" -o /rootfs/init && strip /rootfs/init

# Copy in primary config, and default secondary config to rootfs
ARG CONFIG_PRIMARY=CONFIG_PRIMARY_UNSET
COPY /config/primary/$CONFIG_PRIMARY /rootfs/config/primary.json
COPY /config/secondary/default.json /rootfs/config/secondary.json
RUN tree /rootfs > ${OUT_DIR}/rootfs.txt

# Build initramfs
RUN if [ -f "/initramfs.cpio" ]; then rm /initramfs.cpio; fi
RUN cd /rootfs && find . -print0 | cpio --null --create --verbose --format=newc > /initramfs.cpio

# Build final kernel with real initramfs
ARG CONFIG_COMPRESSION=GZIP
RUN cd linux-${VERSION_KERNEL} && \
    make CONFIG_KERNEL_${CONFIG_COMPRESSION}=y CONFIG_INITRAMFS_COMPRESSION_${CONFIG_COMPRESSION}=y CFLAGS="-pipe -Os -s -fno-stack-protector -U_FORTIFY_SOURCE" KGZIP=pigz -j $(nproc) && \
    cp arch/x86_64/boot/bzImage ${OUT_DIR}/BOOTx64-$CONFIG_MODULES-$CONFIG_PRIMARY-$CONFIG_COMPRESSION.EFI && rm arch/x86_64/boot/bzImage && \
    cd ${OUT_DIR} && ln -s BOOTx64-$CONFIG_MODULES-$CONFIG_PRIMARY-$CONFIG_COMPRESSION.EFI BOOTx64.EFI

USER 100:100
CMD ["cp", "-r" ,"/build/out", "/"]
