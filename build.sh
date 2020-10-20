#!/bin/sh

set -ex

KERNEL_CONFIG=config-5.9-RT-testing

KERNEL_VERSION=5.9  
KERNEL_RT=5.9.1-rt17 
MUSL_VERSION=1.2.0
IPTABLES_VERSION=1.8.5
DOCKER_VERSION=19.03.9
KMOD=26
WG=v1.0.20200513
AUFS=aufs5-standalone

NUM_JOBS="$(grep ^processor /proc/cpuinfo | wc -l)"

BUILD_DIR=/build
ROOTFS_DIR=$BUILD_DIR/rootfs
SRC_DIR=$BUILD_DIR/src
OUT_DIR=$BUILD_DIR/out
AUFS_SRC=$SRC_DIR/$AUFS

debug() {
  echo "Dropping into a shell for debugging ..."
  /bin/sh
}

# download_wg() { 
#   cd $SRC_DIR
#   if [ ! -f "wireguard-tools" ]; then
#     git clone -b $WG https://git.zx2c4.com/wireguard-tools wireguard-tools
#   fi
# }

download_kmod() {
  cd $SRC_DIR
  if [ ! -f "kmod.tar.xz" ]; then
    wget -q -O kmod.tar.xz \
      https://mirrors.edge.kernel.org/pub/linux/utils/kernel/kmod/kmod-$KMOD.tar.xz
    tar -xf kmod.tar.xz
  fi
}

download_kernel() {
  cd $SRC_DIR
  if [ ! -f "kernel.tar.xz" ]; then
    wget -q -O kernel.tar.xz http://kernel.org/pub/linux/kernel/v5.x/linux-$KERNEL_VERSION.tar.xz
    tar -xf kernel.tar.xz
    wget -q -O patch-$KERNEL_RT.patch.xz https://cdn.kernel.org/pub/linux/kernel/projects/rt/5.9/patch-$KERNEL_RT.patch.xz
  
  fi
}

download_musl() {
  cd $SRC_DIR
  if [ ! -f "musl.tar.gz" ]; then
    wget -q -O musl.tar.gz \
      http://www.musl-libc.org/releases/musl-$MUSL_VERSION.tar.gz
    tar -xf musl.tar.gz
  fi
}

download_iptables() {
  if [ ! -f "iptables.tar.bz2" ]; then
    cd $SRC_DIR
    wget -q -O iptables.tar.bz2 \
      https://netfilter.org/projects/iptables/files/iptables-$IPTABLES_VERSION.tar.bz2
    tar -xf iptables.tar.bz2
  fi
}

download_docker() {
  if [ ! -f "docker.tgz" ]; then
    cd $SRC_DIR
    wget -q -O docker.tgz \
      https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION.tgz
    tar -xf docker.tgz
  fi
}

#download_aufs() {
#  cd $SRC_DIR
#  git clone git://github.com/sfjro/aufs5-standalone.git $AUFS
#  cd $AUFS
#  git checkout origin/aufs5.6
#}

# build_wg() {
#   (
#     cd $SRC_DIR/wireguard-tools/src
#     make -j $NUM_JOBS
#     make DESTDIR=$ROOTFS_DIR install
#   )
# }

build_musl() {
  (
    cd $SRC_DIR/musl-$MUSL_VERSION
    ./configure \
      --prefix=/usr
    make -j $NUM_JOBS
    make DESTDIR=$ROOTFS_DIR install
  )
}

build_iptables() {
  (
    cd $SRC_DIR/iptables-$IPTABLES_VERSION

    ./configure \
      --prefix=/usr \
      --enable-libipq \
      --disable-nftables \
      --with-xtlibdir=/lib/xtables

    make \
      EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" \
      -j $NUM_JOBS
    make DESTDIR=$ROOTFS_DIR install
  )
}

install_docker() {
  mkdir -p $ROOTFS_DIR/usr/bin/
  cp $SRC_DIR/docker/* $ROOTFS_DIR/usr/bin/

  strip $ROOTFS_DIR/usr/bin/containerd 
  strip $ROOTFS_DIR/usr/bin/containerd-shim 
  strip $ROOTFS_DIR/usr/bin/ctr 
  strip $ROOTFS_DIR/usr/bin/docker 
  strip $ROOTFS_DIR/usr/bin/docker-init 
  strip $ROOTFS_DIR/usr/bin/docker-proxy 
  strip $ROOTFS_DIR/usr/bin/dockerd
  strip $ROOTFS_DIR/usr/bin/runc
}

build_rootfs() {
  (
    cd $ROOTFS_DIR

    # Cleanup rootfs
    find . -type f -name '.empty' -size 0c -delete
    rm -rf usr/man usr/share/man
    rm -rf usr/lib/pkgconfig
    rm -rf usr/include

    u-root -initcmd="/init-custom" -uinitcmd="/uinit-custom" -build=bb -format=cpio -o /build/initrmfs.cpio -files $ROOTFS_DIR:/ core boot
  )
}

patch_kernel() {
    cd $SRC_DIR/linux-$KERNEL_VERSION
    
    # RT_PREEMPT and AUFS doesn't play nicely togther.

    #echo "AUFS patching  " + $AUFS_SRC

    #cat $AUFS_SRC/aufs5-base.patch | patch -Np1
    #cat $AUFS_SRC/aufs5-kbuild.patch | patch -Np1
    #cat $AUFS_SRC/aufs5-mmap.patch | patch -Np1
    #cat $AUFS_SRC/aufs5-standalone.patch | patch -Np1

    #rm -f $AUFS_SRC/include/uapi/linux/Kbuild

    #cp -av $AUFS_SRC/Documentation Documentation/
    #cp -av $AUFS_SRC/fs/* fs/
    #cp -av $AUFS_SRC/include/* include/

    xzcat ../patch-$KERNEL_RT.patch.xz | patch -p1

}



build_kernel() {
  (
    cd $SRC_DIR/linux-$KERNEL_VERSION

    cp $BUILD_DIR/$KERNEL_CONFIG .config


    #make mrproper defconfig -j $NUM_JOBS
    # NOT NEEDED WITH IF THE KERNEL CONFIG IS CORRECRTLY CONFIGURED
    make oldconfig -j $NUM_JOBS

    # finally build the kernel
    make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS
    make INSTALL_MOD_PATH=$ROOTFS_DIR modules_install
    # create the initrmfs

    #cp `find /build/rootfs/ -name e1000e.ko` $rootfs/lib/modules
    #cp `find /build/rootfs/ -name e1000.ko` $rootfs/lib/modules
    #cp `find /build/rootfs/ -name btrfs.ko` $rootfs/lib/modules
    #cp `find /build/rootfs/ -name hid-generic.ko` $rootfs/lib/modules
    #cp `find /build/rootfs/ -name input-leds.ko` $rootfs/lib/modules
    #cp `find /build/rootfs/ -name igb.ko` $rootfs/lib/modules

    #rm -rf /build/rootfs/lib/modules/4.20.12-mjolnir

    u-root -initcmd="/init-custom" -uinitcmd="/uinit-custom" -build=bb -format=cpio -o /build/initrmfs.cpio -files $ROOTFS_DIR:/ core boot

    make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS

#    cp arch/x86_64/boot/bzImage $SRC_DIR/kernel.gz
    cp arch/x86_64/boot/bzImage $OUT_DIR/BOOTx64.EFI
  )
}


rebuild_kernel() {
    cd $SRC_DIR/linux-$KERNEL_VERSION
    u-root -initcmd="/init-custom" -uinitcmd="/uinit-custom" -build=bb -format=cpio -o /build/initrmfs.cpio -files $ROOTFS_DIR:/ core boot

    make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS

#    cp arch/x86_64/boot/bzImage $SRC_DIR/kernel.gz
    cp arch/x86_64/boot/bzImage $OUT_DIR/BOOTx64.EFI
}

build_custom_init() {
  (
    cd $BUILD_DIR/init/init-custom
    go get github.com/u-root/u-root
    go build -o $ROOTFS_DIR/init-custom
    
    cd $BUILD_DIR/init/uinit-custom
    go build -o $ROOTFS_DIR/uinit-custom

    strip $ROOTFS_DIR/uinit-custom
    strip $ROOTFS_DIR/init-custom
  )
}

build_kmod() {
  (
    cd $SRC_DIR/kmod-$KMOD
    ./configure --prefix=/usr --bindir=/bin --sysconfdir=/etc --with-rootlibdir=/lib
    make EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE"
    make DESTDIR=$ROOTFS_DIR install

    cd $ROOTFS_DIR
    for target in depmod insmod lsmod modinfo modprobe rmmod; do
      ln -sfv ../bin/kmod sbin/$target
    done

  )
}

download_packages() {
  download_musl
  download_kmod
  download_iptables
  download_kernel
  download_docker
 # download_aufs
  install_docker
}

build_packages() {
  build_musl
  build_iptables
  build_kmod
}

clean() {
  [ -d "/build/src" ] || rm -rf "/build/src"
}

prepare_build() {
  clean
  mkdir -p $SRC_DIR

  # Clean up old out
  rm -rf $OUT_DIR/*

  # Clean up old rootfs
  #mkdir -p $ROOTFS_DIR
  #rm -rf $ROOTFS_DIR/*
}

rebuild_system(){
  # build the Golang init command
  build_custom_init

  rebuild_kernel
}

build_all() {

  prepare_build

  download_packages
  patch_kernel

  build_packages

  # build the Golang init command
  build_custom_init

  # Creates the release file in the rootfs
  # makes the rootfs into an initramfs, which is build into the kernel (see kernel config)
  build_rootfs

  # makes the kernel into EFI image (/img/BOOTx64.EFI) which can deployed directly on a target system on a VFAT EFI partition in the location /EFI/BOOT/BOOTx64.EFI
  build_kernel
}

case "${1}" in
prepare)
  prepare_build
  ;;
download)
  download_packages
  ;;
patch)
  patch_kernel
  ;;
build)
  build_packages
  ;;
build_kernel)
  build_kernel
  ;;
build_init)
  build_custom_init
  ;;  
rebuild)
  rebuild_system
  ;;
*)
  build_all
  ;;
esac