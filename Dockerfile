# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: Apache-2.0

FROM alpine:3.14.0

RUN apk --no-cache add \
    alpine-sdk argp-standalone asciidoc autoconf automake bc bison build-base ccache clang cmake cryptsetup coreutils \
    device-mapper-static diffutils docbook2x e2fsprogs elfutils-dev flex gawk gettext-dev git gmp-dev go gnupg \
    gzip json-c json-c-dev libaio-dev libmnl-dev libnfnetlink-dev libnftnl-dev libpciaccess-dev libtool \
    libuuid linux-headers llvm llvm-dev lld lvm2-dev lvm2-static lz4 lzo ncurses-dev openssl openssl-dev openssl-libs-static \
    perl pzstd pigz popt popt popt-dev readline-dev rsync tpm2-tss tpm2-tss-dev tpm2-tss-esys tpm2-tss-fapi tpm2-tss-mu \
    tpm2-tss-sys tree upx util-linux util-linux-dev wget xorriso xz zstd 

SHELL ["/bin/bash", "-c"]
RUN git config --global advice.detachedHead false
ENV CC=/usr/bin/gcc CXX=/usr/bin/g++

# Dirs
ENV SRC_DIR=/build/src
ENV OUT_DIR=/build/out
RUN mkdir -p ${OUT_DIR} && mkdir -p ${SRC_DIR} && mkdir -p /rootfs
WORKDIR ${SRC_DIR}

# Package versions
ENV VERSION_KERNEL=5.10.41
ENV VERSION_RT=5.10.41-rt42
ENV VERSION_MUSL=1.2.2
ENV VERSION_DOCKER=20.10.7
ENV VERSION_BUSYBOX=1.33.1
ENV VERSION_WGTOOLS=v1.0.20210424
ENV VERSION_MICROCODE_INTEL=20210608
ENV VERSION_IPTABLES=1.8.7

# Flags
ENV CONFIG_KERNEL=5.10.1-rt20

# Download sources
RUN wget -q -O kernel.tar.xz https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-${VERSION_KERNEL}.tar.xz
RUN wget -q -O patch-rt.xz https://cdn.kernel.org/pub/linux/kernel/projects/rt/5.10/older/patch-${VERSION_RT}.patch.xz
RUN wget -q -O musl.tar.gz https://www.musl-libc.org/releases/musl-${VERSION_MUSL}.tar.gz
RUN wget -q -O docker.tgz https://download.docker.com/linux/static/stable/x86_64/docker-${VERSION_DOCKER}.tgz
RUN wget -q -O wireguard.tar.xz https://git.zx2c4.com/wireguard-tools/snapshot/wireguard-tools-${VERSION_WGTOOLS}.tar.xz
RUN wget -q -O iptables.tar.bz2  https://netfilter.org/projects/iptables/files/iptables-${VERSION_IPTABLES}.tar.bz2
RUN wget -q -O busybox.tar.bz2 https://busybox.net/downloads/busybox-${VERSION_BUSYBOX}.tar.bz2
RUN wget -q -O microcode.tar.gz https://github.com/intel/Intel-Linux-Processor-Microcode-Data-Files/archive/refs/tags/microcode-${VERSION_MICROCODE_INTEL}.tar.gz 

# Verify sources
RUN echo "f604759de80767c4f8bdc500eec730dc161bc914a48bd366b748c176701a6771 kernel.tar.xz" | sha256sum -c -
RUN echo "03a1be966680c3fc8853d8b1d08fca3dd1303961e471d5bb41e44d57b07e12fd patch-rt.xz" | sha256sum -c -
RUN echo "9b969322012d796dc23dda27a35866034fa67d8fb67e0e2c45c913c3d43219dd musl.tar.gz" | sha256sum -c -
RUN echo "34ad50146fce29b28e5115a1e8510dd5232459c9a4a9f28f65909f92cca314d9 docker.tgz" | sha256sum -c -
RUN echo "98140aa91ea04018ebd874c14ab9b6994f48cdaf9a219ccf7c0cd3e513c7428a wireguard.tar.xz" | sha256sum -c -
RUN echo "c109c96bb04998cd44156622d36f8e04b140701ec60531a10668cfdff5e8d8f0 iptables.tar.bz2" | sha256sum -c -
RUN echo "12cec6bd2b16d8a9446dd16130f2b92982f1819f6e1c5f5887b6db03f5660d28 busybox.tar.bz2" | sha256sum -c -
RUN echo "fd85b6b769efd029dec6a2c07106fd18fb4dcb548b7bc4cde09295a8344ef6d7 microcode.tar.gz" | sha256sum -c -

# Patch kernel
RUN tar -xf kernel.tar.xz
RUN cd linux-${VERSION_KERNEL} && xzcat ../patch-$KERNEL_RT.patch.xz | patch -p1

# Build kernel with custom config
RUN touch /initramfs.cpio
COPY config/config-${CONFIG_KERNEL} linux-${VERSION_KERNEL}/.config
RUN cd linux-${VERSION_KERNEL} && make oldconfig -j $(nproc)
RUN cd linux-${VERSION_KERNEL} && make CFLAGS="-pipe -Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $(nproc)

# Build musl
RUN tar -xf musl.tar.gz
RUN cd musl-${VERSION_MUSL} && ./configure --prefix=/usr
RUN cd musl-${VERSION_MUSL} && make -j $(nproc)

# Build docker
RUN tar -xf docker.tgz

# Build wg-tools
RUN tar -xf wireguard.tar.xz
RUN cd wireguard-tools-${VERSION_WGTOOLS}/src && make -j $(nproc)

# Build busybox
RUN tar -xf busybox.tar.bz2
COPY config/busybox-config busybox-${VERSION_BUSYBOX}/.config
RUN cd busybox-${VERSION_BUSYBOX} && make oldconfig
RUN cd busybox-${VERSION_BUSYBOX} && make -j $(nproc)

# Build iptables
RUN tar -xf iptables.tar.bz2
RUN cd iptables-${VERSION_IPTABLES} && ./configure --prefix=/ --mandir=/tmp --disable-nftables
RUN cd iptables-${VERSION_IPTABLES} && make EXTRA_CFLAGS="-pipe -Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $(nproc)

# Set up template rootfs
COPY template_rootfs /template_rootfs
COPY scripts/build-rootfs.sh .
RUN ./build-rootfs.sh /rootfs /template_rootfs

# Install CPU microcode to rootfs
RUN tar -xzf microcode.tar.gz
RUN mkdir -p /lib/firmware/intel-ucode && cp -r Intel-Linux-Processor-Microcode-Data-Files-microcode-${VERSION_MICROCODE_INTEL}/intel-ucode /lib/firmware/intel-ucode

# Install modules to rootfs
RUN cd linux-${VERSION_KERNEL} && mkdir -p /rootfs && make -j $(nproc) INSTALL_MOD_PATH=/rootfs modules_install

# Install packages to rootfs
RUN cd busybox-${VERSION_BUSYBOX} && make -j $(nproc) CONFIG_PREFIX=/rootfs install
RUN cd musl-${VERSION_MUSL} && make -j $(nproc) DESTDIR=/rootfs install
RUN cd iptables-${VERSION_IPTABLES} && make DESTDIR=/rootfs install
RUN cp wireguard-tools-${VERSION_WGTOOLS}/src/wg /rootfs/usr/sbin/wg
RUN cp docker/* /rootfs/usr/bin/

# Add alpine packages
RUN cd /bin && cp -t /rootfs/bin lsblk 
RUN cd /lib && cp -t /rootfs/lib libblkid.so.* libsmartcols.so.* libmount.so.*

# Strip modules if specified
ARG CONFIG_MODULES=ALL
COPY config/modules .
RUN find /rootfs/lib/modules | grep "\.ko$" > ${OUT_DIR}/modules.txt
RUN if [ "${CONFIG_MODULES}" != "ALL" ]; then find /rootfs/lib/modules | grep "\.ko$" | grep -v -f ${CONFIG_MODULES} | xargs rm; fi;
RUN find /rootfs/lib/modules | grep "\.ko$" > ${OUT_DIR}/modules_selected.txt
RUN find /rootfs/lib/modules | grep "\.ko$" | xargs du -sh | sort -rh > ${OUT_DIR}/modules_selected_size.txt

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
COPY config/primary/$CONFIG_PRIMARY /rootfs/config/primary.yml
COPY config/secondary/default.yml /rootfs/config/secondary.yml
RUN find /rootfs > ${OUT_DIR}/rootfs.txt

# Build initramfs
RUN if [ -f "/initramfs.cpio" ]; then rm /initramfs.cpio; fi
RUN cd /rootfs && find . -print0 | cpio --null --create --verbose --format=newc > /initramfs.cpio

# Optionally patch compression ratio to speed up build
ARG COMPRESSION_LEVEL="22"
RUN cd linux-${VERSION_KERNEL} && sed -i -e 's/$(ZSTD) -22 --ultra/$(ZSTD) -T0 -${COMPRESSION_LEVEL} --ultra/g' scripts/Makefile.lib 

# Build final kernel with real initramfs
RUN cd linux-${VERSION_KERNEL} && make CFLAGS="-pipe -Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $(nproc) 

# Create BOOTx64.EFI symlink. MD5 is fine as this is not for security, only for a 128 bit hash filename
RUN HASH=($(md5sum linux-${VERSION_KERNEL}/arch/x86_64/boot/bzImage)) && cp linux-${VERSION_KERNEL}/arch/x86_64/boot/bzImage ${OUT_DIR}/${HASH}.EFI && \
    cd ${OUT_DIR} && ln -s ${HASH}.EFI BOOTx64.EFI && ln -s ${HASH}.EFI BOOTx64-$CONFIG_MODULES-$CONFIG_PRIMARY.ZSTD$COMPRESSION_LEVEL.EFI

FROM alpine:3.14.0
COPY --from=0 /build/out /build/out

USER 100:100
CMD ["cp", "-r" ,"/build/out", "/"]