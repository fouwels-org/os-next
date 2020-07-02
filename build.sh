#!/bin/sh

set -ex

KERNEL_VERSION=5.2.14
MUSL_VERSION=1.2.0
IPTABLES_VERSION=1.8.5
DOCKER_VERSION=19.03.9
KMOD=26

NUM_JOBS="$(grep ^processor /proc/cpuinfo | wc -l)"

build=/build
rootfs=$build/rootfs/
isoimage=$build/isoimage/

debug() { echo "Dropping into a shell for debugging ..."; /bin/sh; }

download_kmod() {
  cd /build/src
  wget -q -O kmod.tar.xz \
    https://mirrors.edge.kernel.org/pub/linux/utils/kernel/kmod/kmod-$KMOD.tar.xz
  tar -xf kmod.tar.xz
}

download_kernel() {
  cd /build/src
  wget -q -O kernel.tar.xz \
    http://kernel.org/pub/linux/kernel/v5.x/linux-$KERNEL_VERSION.tar.xz
  tar -xf kernel.tar.xz
}

download_musl() {
  cd /build/src    
  wget -q -O musl.tar.gz \
    http://www.musl-libc.org/releases/musl-$MUSL_VERSION.tar.gz
  tar -xf musl.tar.gz
}

download_iptables() {
  cd /build/src    
  wget -q -O iptables.tar.bz2 \
    https://netfilter.org/projects/iptables/files/iptables-$IPTABLES_VERSION.tar.bz2
  tar -xf iptables.tar.bz2
}

download_docker() {
  cd /build/src    
  wget -q -O docker.tgz \
    https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION.tgz
  tar -xf docker.tgz
}

create_rootfs() {
  (
  test -d "$rootfs" || mkdir "$rootfs"
  cd $rootfs
  mkdir -p dev/pts
  mkdir -p lib
  mkdir -p etc/dropbear
  mkdir -p etc/ssl/certs
  mkdir -p mnt
  mkdir -p etc
  mkdir -p proc
  mkdir -p run
  mkdir -p sys
  mkdir -p sbin
  mkdir -p var/lib/docker
  mkdir -p var/log
  mkdir -p var/run
  install -d -m 0750 root
  install -d -m 1777 tmp
  mkdir -p usr/bin
  mkdir -p usr/lib
  mkdir -p usr/share/udhcpc
  
  cd /build
  echo rootfs/dev/pts rootfs/etc/dropbear rootfs/lib rootfs/mnt rootfs/proc rootfs/root rootfs/run rootfs/sys rootfs/tmp rootfs/usr/lib rootfs/var/lib/docker   | xargs -n 1 cp /build/files/.empty
  cd /build/files
  cp fstab group issue mime.types motd mtab passwd profile rc.conf securetty shadow shells $rootfs/etc
  cp lastlog /build/rootfs/var/log
  cp utmp /build/rootfs/var/run
  cp mail /build/rootfs/var

  cp etc/* /build/rootfs/etc/

  cp ca-bundle.pem cert.pem /build/rootfs/etc/ssl/certs
  )
}

build_musl() {
  (
  cd musl-$MUSL_VERSION
  ./configure \
    --prefix=/usr
  make -j $NUM_JOBS
  make DESTDIR=$rootfs install
  )
}

build_iptables() {
  (
  cd iptables-$IPTABLES_VERSION
  ./configure  \
    --prefix=/usr \
    --enable-libipq \
    --disable-nftables \
    --with-xtlibdir=/lib/xtables

  make \
    EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" \
    -j $NUM_JOBS
  make DESTDIR=$rootfs install
  )
}

install_docker() {
  mv docker/* $rootfs/usr/bin/
}

build_rootfs() {
  (
   
  cd $rootfs

  # Cleanup rootfs
  find . -type f -name '.empty' -size 0c -delete
  rm -rf usr/man usr/share/man
  rm -rf usr/lib/pkgconfig
  rm -rf usr/include
  u-root -initcmd="/uinit-custom" -build=bb -format=cpio -o /build/initrmfs.cpio -files $rootfs:/ core boot

  )
}

build_kernel() {
  (
  cd /build/src
  cd linux-$KERNEL_VERSION
  cp /build/config .config
  #make mrproper defconfig -j $NUM_JOBS
  # NOT NEEDED WITH IF THE KERNEL CONFIG IS CORRECRTLY CONFIGURED
  make oldconfig -j $NUM_JOBS

  # finally build the kernel
  make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS
  make INSTALL_MOD_PATH=$rootfs modules_install
  # create the initrmfs
  u-root -uinitcmd="/uinit-custom" -build=bb -format=cpio -o /build/initrmfs.cpio -files $rootfs:/ core boot 

  make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS

  cp arch/x86_64/boot/bzImage /build/src/kernel.gz
  cp arch/x86_64/boot/bzImage /img/BOOTx64.EFI
  )
}

build_custom_init(){
  (
  cd /build/uinit-custom
  go build 
  cp uinit-custom $rootfs/uinit-custom
  )
}

build_kmod() {
  (
  cd kmod-$KMOD
  ./configure --prefix=/usr --bindir=/bin --sysconfdir=/etc --with-rootlibdir=/lib
  make EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE"
  make DESTDIR=$rootfs install
  cd /build/rootfs
  for target in depmod insmod lsmod modinfo modprobe rmmod; do
    ln -sfv ../bin/kmod sbin/$target
  done

  )
}


download_packages(){
  cd /build/src
  download_musl
  download_kmod
  download_iptables
  download_kernel
  download_docker
  install_docker
}

build_packages(){
  cd /build/src
  build_musl
  build_iptables
  build_kmod
}

clean(){
  [ -d $rootfs ] || rm -rf $rootfs
  [ -d "/build/src" ] || rm -rf "/build/src"
}

prepare_build(){
  clean
  mkdir /build/src
  # Make sure the rootfs is clean
  create_rootfs
  cd /build/src
}

build_all() {
  # WORDIR build/src
  prepare_build
  
  # WORDIR build/src
  download_packages

  # WORDIR build/src
  build_packages

  # build the Golang init command
  build_custom_init
  # WORDIR $rootfs
  # Creates the release file in the rootfs
  # makes the rootfs into an initramfs, which is build into the kernel (see kernel config)
  build_rootfs

  # WORDIR build/src
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
  build)
    build_packages
    ;;
  build_kernel)
    build_kernel
    ;;
  repack)
    repack
    ;;
  *)
    build_all
    ;;
esac
