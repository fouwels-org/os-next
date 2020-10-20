#!/bin/sh

set -ex
PS4="[main] "

KERNEL_CONFIG=config-5.9-RT-testing

KERNEL_VERSION=5.9  
KERNEL_RT=5.9.1-rt17 
MUSL_VERSION=1.2.0
DOCKER_VERSION=19.03.9
KMOD=26
AUFS=aufs5-standalone

NFTABLES_TAG=v0.9.6
NFTABLES_LIBNFTNL_TAG=libnftnl-1.1.7
NFTABLES_LIBMNL_TAG=libmnl-1.0.4
NFTABLES_IPTABLES_VERSION=v1.8.5


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

download_nftables() {
  # surpress git detached head message, this is what we want
  git config --global advice.detachedHead false

  if [ ! -d "$SRC_DIR/nftables" ]; then
    cd $SRC_DIR
    git clone --depth 1 --branch $NFTABLES_TAG git://git.netfilter.org/nftables
    git clone --depth 1 --branch $NFTABLES_LIBNFTNL_TAG git://git.netfilter.org/libnftnl
    git clone --depth 1 --branch $NFTABLES_LIBMNL_TAG git://git.netfilter.org/libmnl
    git clone --depth 1 --branch $NFTABLES_IPTABLES_VERSION git://git.netfilter.org/iptables
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

build_musl() {
    cd $SRC_DIR/musl-$MUSL_VERSION
    ./configure \
      --prefix=/usr
    make -j $NUM_JOBS
    make DESTDIR=$ROOTFS_DIR install
}

build_kmod() {
  
  cd $SRC_DIR/kmod-$KMOD
  ./configure --prefix=/usr --bindir=/bin --sysconfdir=/etc --with-rootlibdir=/lib
  make EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE"
  make DESTDIR=$ROOTFS_DIR install

  cd $ROOTFS_DIR
  for target in depmod insmod lsmod modinfo modprobe rmmod; do
    ln -sfv ../bin/kmod sbin/$target
  done
}

build_nftables() {
  cd $SRC_DIR/libnftnl && sh autogen.sh && ./configure && make -j $NUM_JOBS && make install
  cd $SRC_DIR/libmnl && sh autogen.sh && ./configure && make -j $NUM_JOBS && make install
  cd $SRC_DIR/nftables && sh autogen.sh && ./configure && make -j $NUM_JOBS && make install
  cd $SRC_DIR/iptables && sh autogen.sh && ./configure && make -j $NUM_JOBS && make install

  # symlink iptables to iptables-nft (nft backed), instead of iptables-legacy (iptables backed)
  # see: https://www.redhat.com/en/blog/using-iptables-nft-hybrid-linux-firewall
  # this will allow docker to call legacy iptables, and write into the nft instead.

  PREFIX=/usr/local/sbin

  rm $PREFIX/iptables
  rm $PREFIX/iptables-save
  rm $PREFIX/iptables-restore
  rm $PREFIX/ip6tables
  rm $PREFIX/ip6tables-save
  rm $PREFIX/ip6tables-restore
  rm $PREFIX/arptables
  rm $PREFIX/arptables-save
  rm $PREFIX/arptables-restore
  rm $PREFIX/ebtables
  rm $PREFIX/ebtables-save
  rm $PREFIX/ebtables-restore

  ln -s $PREFIX/iptables-nft $PREFIX/iptables
  ln -s $PREFIX/iptables-nft-save $PREFIX/iptables-save
  ln -s $PREFIX/iptables-nft-restore $PREFIX/iptables-restore
  ln -s $PREFIX/ip6tables-nft $PREFIX/ip6tables
  ln -s $PREFIX/ip6tables-nft-save $PREFIX/ip6tables-save
  ln -s $PREFIX/ip6tables-nft-restore $PREFIX/ip6tables-restore
  ln -s $PREFIX/arptables-nft $PREFIX/arptables
  ln -s $PREFIX/arptables-nft-save $PREFIX/arptables-save
  ln -s $PREFIX/arptables-nft-restore $PREFIX/arptables-restore
  ln -s $PREFIX/ebtables-nft $PREFIX/ebtables
  ln -s $PREFIX/ebtables-nft-save $PREFIX/ebtables-save
  ln -s $PREFIX/ebtables-nft-restore $PREFIX/ebtables-restore
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
  PS4="[build_rootfs] "

  cd $ROOTFS_DIR

  # Cleanup rootfs
  find . -type f -name '.empty' -size 0c -delete
  rm -rf usr/man usr/share/man
  rm -rf usr/lib/pkgconfig
  rm -rf usr/include

  u-root -initcmd="/init-custom" -uinitcmd="/uinit-custom" -build=bb -format=cpio -o /build/initrmfs.cpio -files $ROOTFS_DIR:/ core boot

  touch $SRC_DIR/flag_build_rootfs
}

patch_kernel() {
  PS4="[patch_kernel] "

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

  touch $SRC_DIR/flag_patched_kernel
}

build_modules() {
  PS4="[build_modules] "
  
  cd $SRC_DIR/linux-$KERNEL_VERSION
  cp -f $BUILD_DIR/$KERNEL_CONFIG .config

  #make mrproper defconfig -j $NUM_JOBS
  # NOT NEEDED WITH IF THE KERNEL CONFIG IS CORRECRTLY CONFIGURED
  make oldconfig -j $NUM_JOBS

  # finally build the kernel
  make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS
  make INSTALL_MOD_PATH=$ROOTFS_DIR modules_install

  touch $SRC_DIR/flag_built_modules
}

build_kernel() {
  PS4="[build_kernel] "

  cd $SRC_DIR/linux-$KERNEL_VERSION
  cp -f $BUILD_DIR/$KERNEL_CONFIG .config

  u-root -initcmd="/init-custom" -uinitcmd="/uinit-custom" -build=bb -format=cpio -o /build/initrmfs.cpio -files $ROOTFS_DIR:/ core boot

  make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS

  #cp arch/x86_64/boot/bzImage $SRC_DIR/kernel.gz
  cp arch/x86_64/boot/bzImage $OUT_DIR/BOOTx64.EFI
}

build_custom_init() {
  PS4="[build_custom_init] "

  cd $BUILD_DIR/init/init-custom
  go get github.com/u-root/u-root
  go build -o $ROOTFS_DIR/init-custom
  
  cd $BUILD_DIR/init/uinit-custom
  go build -o $ROOTFS_DIR/uinit-custom

  strip $ROOTFS_DIR/uinit-custom
  strip $ROOTFS_DIR/init-custom
}

download_packages() {
  PS4="[download_packages] "

  download_musl
  download_kmod
  download_nftables
  download_kernel
  download_docker
  #download_aufs
  install_docker

  touch $SRC_DIR/flag_downloaded
}

build_packages() {
  PS4="[build_packages] "

  build_musl
  build_kmod
  build_nftables

  touch $SRC_DIR/flag_built_packages
}

clean() {
  [ -d "/build/src" ] || rm -rf "/build/src"
}

prepare_build() {
  PS4="[prepare_build] "
  # Create src dir
  if [ ! -d $SRC_DIR ]; then
    mkdir -p $SRC_DIR
  fi  

  # Clean up old out
  if [ ! -d $OUT_DIR ]; then
    mkdir -p $OUT_DIR
  fi
  rm -rf $OUT_DIR/*

  # Clean up old rootfs
  #mkdir -p $ROOTFS_DIR
  #rm -rf $ROOTFS_DIR/*
}

build_all() {

  prepare_build

  if [ ! -f $SRC_DIR/flag_downloaded ]; then
    download_packages
  fi
  
  if [ ! -f $SRC_DIR/flag_patched_kernel ]; then
    patch_kernel
  fi

  if [ ! -f $SRC_DIR/flag_built_packages ]; then
    build_packages
  fi

  # build the Golang init command
  build_custom_init

  # Creates the release file in the rootfs
  # makes the rootfs into an initramfs, which is build into the kernel (see kernel config)
  if [ ! -f $SRC_DIR/flag_built_rootfs ]; then
    build_rootfs
  fi

  if [ ! -f $SRC_DIR/flag_built_modules ]; then
    build_modules
  fi

  # makes the kernel into EFI image (/img/BOOTx64.EFI) which can deployed directly on a target system on a VFAT EFI partition in the location /EFI/BOOT/BOOTx64.EFI
  build_kernel
}

# set -e needs to be re-applied after every ), as the bracket creates a new scope.
case "${1}" in
prepare)
  set -ex
  prepare_build
  ;;
download)
  set -ex
  download_packages
  ;;
patch)
  set -ex
  patch_kernel
  ;;
build_packages)
  set -ex
  build_packages
  ;;
build_kernel)
  set -ex
  build_modules
  build_kernel
  ;;
build_init)
  set -ex
  build_custom_init
  ;;
*)
  set -ex
  build_all
  ;;
esac