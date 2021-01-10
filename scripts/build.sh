#!/bin/sh

set -ex
PS4="[main] "

KERNEL_CONFIG=config-5.10.1-rt20

# Kernel versions
KERNEL_VERSION=5.10.1
KERNEL_RT=5.10.1-rt20

# base system versions
MUSL_VERSION=1.2.1
DOCKER_VERSION=20.10.2
BUSYBOX_VERSION=1.32.1

# security system versions, full disk encryption and wireguard P2P
WG_TOOLS=v1.0.20200827

# Networking for docker (nftables rather than IPtables), using kernel nftables
NFTABLES_TAG=v0.9.6
NFTABLES_LIBNFTNL_TAG=libnftnl-1.1.7
NFTABLES_LIBMNL_TAG=libmnl-1.0.4
NFTABLES_IPTABLES_VERSION=v1.8.5

NUM_JOBS="$(grep ^processor /proc/cpuinfo | wc -l)"

BUILD_DIR=/build
ROOTFS_DIR=$BUILD_DIR/rootfs
SRC_DIR=$BUILD_DIR/src
OUT_DIR=$BUILD_DIR/out

GOPATH=$SRC_DIR/go
GOBIN=$GOPATH/bin
PATH=$PATH:$GOBIN

debug() {
  echo "Dropping into a shell for debugging ..."
  /bin/sh
}

###############################################
# Download source and packages
###############################################

download_busybox() {
  cd /build/src    
  wget -q -O busybox.tar.bz2 \
    http://busybox.net/downloads/busybox-$BUSYBOX_VERSION.tar.bz2
  tar -xf busybox.tar.bz2
}

download_kernel() {
  cd $SRC_DIR
  if [ ! -f "kernel.tar.xz" ]; then
    wget -q -O kernel.tar.xz http://kernel.org/pub/linux/kernel/v5.x/linux-$KERNEL_VERSION.tar.xz
    tar -xf kernel.tar.xz
    
    wget -q -O patch-$KERNEL_RT.patch.xz https://cdn.kernel.org/pub/linux/kernel/projects/rt/5.10/older/patch-$KERNEL_RT.patch.xz
    
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

  if [ ! -d "$SRC_DIR/iptables" ]; then
    cd $SRC_DIR
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

download_wg_tools(){
  if [ ! -f "wireguard-tools.tar.xz" ]; then
    cd $SRC_DIR
    git clone -b $WG_TOOLS https://git.zx2c4.com/wireguard-tools
  fi
}

###############################################
# Build Packages
###############################################

build_musl() {
    cd $SRC_DIR/musl-$MUSL_VERSION
    ./configure \
      --prefix=/usr
    make -j $NUM_JOBS
}

build_nftables() {
  cd $SRC_DIR/libnftnl && sh autogen.sh && ./configure --prefix=/usr && make -j $NUM_JOBS
  cd $SRC_DIR/libmnl && sh autogen.sh && ./configure --prefix=/usr && make -j $NUM_JOBS
  cd $SRC_DIR/iptables && sh autogen.sh && ./configure --prefix=/usr && make -j $NUM_JOBS 
}

build_wg_tools() {
    cd $SRC_DIR/wireguard-tools/src
    make -j $NUM_JOBS
}

###############################################
# Install Packages
###############################################

install_musl() {
    cd $SRC_DIR/musl-$MUSL_VERSION
    make DESTDIR=$ROOTFS_DIR install
}

install_nftables() {
  cd $SRC_DIR/libnftnl && make DESTDIR=$ROOTFS_DIR install  
  cd $SRC_DIR/libmnl && make DESTDIR=$ROOTFS_DIR install
  cd $SRC_DIR/iptables && make DESTDIR=$ROOTFS_DIR install
  # symlink iptables to iptables-nft (nft backed), instead of iptables-legacy (iptables backed)
  # see: https://www.redhat.com/en/blog/using-iptables-nft-hybrid-linux-firewall
  # this will allow docker to call legacy iptables, and write into the nft instead.

  cd $ROOTFS_DIR
  PREFIX=usr/sbin

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

  ln -sfv ../$PREFIX/iptables-nft sbin/iptables
  ln -sfv ../$PREFIX/iptables-nft-save sbin/iptables-save
  ln -sfv ../$PREFIX/iptables-nft-restore sbin/iptables-restore
  ln -sfv ../$PREFIX/ip6tables-nft sbin/ip6tables
  ln -sfv ../$PREFIX/ip6tables-nft-save sbin/ip6tables-save
  ln -sfv ../$PREFIX/ip6tables-nft-restore sbin/ip6tables-restore
  ln -sfv ../$PREFIX/arptables-nft sbin/arptables
  ln -sfv ../$PREFIX/arptables-nft-save sbin/arptables-save
  ln -sfv ../$PREFIX/arptables-nft-restore sbin/arptables-restore
  ln -sfv ../$PREFIX/ebtables-nft sbin/ebtables
  ln -sfv ../$PREFIX/ebtables-nft-save sbin/ebtables-save
  ln -sfv ../$PREFIX/ebtables-nft-restore sbin/ebtables-restore
}

install_wg_tools() {
    cd $SRC_DIR/wireguard-tools/src
    strip wg
    if [ ! -d $ROOTFS_DIR/usr/sbin ]; then
      mkdir -p $ROOTFS_DIR/usr/sbin
    fi
    cp wg $ROOTFS_DIR/usr/sbin/wg
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

###############################################
# Build Rootfs 
###############################################

build_rootfs() {
  PS4="[build_rootfs] "

  cd $ROOTFS_DIR

  cp -a /build/template_rootfs/* .  

  # template directories contain a file called .empty, to ensure git doesn't ignore the empty directory
  find . -type f -name '.empty' -size 0c -delete
  
  chmod 0755 bin dev etc proc sbin sys usr
  chmod -R 0777 tmp var
  chmod 0770 root

  # create temp character devices to allow for inital boot
  mknod dev/console c 5 1
  chmod 0600 dev/console 
  mknod dev/tty c 5 0
  chmod 0666 dev/tty 
  mknod dev/null c 1 3
  chmod 0666 dev/null 
  mknod dev/port c 1 4
  chmod 0640 dev/port 
  mknod dev/urandom c 1 9
  chmod 0640 dev/urandom 

  # add timezone GMT default
  # This is the literal timezone file for GMT-0. Given that we have no idea where we will be running, GMT seems a reasonable guess. If it
	# matters, setup code should download and change this to something else.
	# GMT0="TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\xf8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00\nGMT0\n"
  #ßecho -n -e $GMT0 > etc/localtime
  chmod 644 etc/localtime

  # Default nameserver set to google. 
  NAMESERVER="nameserver 8.8.8.8"
  echo $NAMESERVER > etc/resolv.conf
  chmod 644 etc/resolv.conf

  #go get github.com/u-root/u-root
  #u-root -initcmd="/init-custom" -defaultsh="" -format=cpio -o /build/initrmfs.cpio -files $ROOTFS_DIR:/
  #cpio -idv < /build/initrmfs.cpio
  chmod 770 $ROOTFS_DIR/usr/share/udhcpc/default.script
  chmod 770 $ROOTFS_DIR/init-custom
  chmod -R 660 $ROOTFS_DIR/etc
  
  ln -sfv init-custom init
  
  # create a fake initrmfs so ensure the kernel will compile the modules before the final initrmfs is created 
  touch /build/initrmfs.cpio

  touch $SRC_DIR/flag_built_rootfs
}

###############################################
# Patch Kernel 
###############################################

patch_kernel() {
  PS4="[patch_kernel] "

  cd $SRC_DIR/linux-$KERNEL_VERSION
  
  # RT_PREEMPT and AUFS doesn't play nicely togther.
  #echo "AUFS patching  " + $AUFS_SRC

  # Adding realtime preempt to the kernel
  xzcat ../patch-$KERNEL_RT.patch.xz | patch -p1

  # Adding wireguard modules to the kernel
  #$SRC_DIR/wireguard-linux-compat/kernel-tree-scripts/create-patch.sh | patch -p1

  touch $SRC_DIR/flag_patched_kernel
}

###############################################
# Build Kernel Modules
###############################################

build_modules() {
  PS4="[build_modules] "
  # create an empty initrmfs.cpio file to trick the kernel build whilst the modules are being build
  touch $BUILD_DIR/initrmfs.cpio

  cd $SRC_DIR/linux-$KERNEL_VERSION
  cp -f $BUILD_DIR/config/$KERNEL_CONFIG .config
  # NOT NEEDED WITH IF THE KERNEL CONFIG IS CORRECRTLY CONFIGURED
  make oldconfig -j $NUM_JOBS
  # finally build the kernel
  make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS
  touch $SRC_DIR/flag_built_modules
}

###############################################
# Install Kernel Modules
###############################################

install_modules() {
  PS4="[install_modules] "
  
  cd $SRC_DIR/linux-$KERNEL_VERSION
  make INSTALL_MOD_PATH=$ROOTFS_DIR modules_install
  # get a system specific modules list and remove all others taht aren't loaded - This is for a production solution
  case "$BUILD_TYPE" in
    "FACTORY")
      echo "including the Factory EFI Image - All modules are included"
      ;;
    "K300")
      echo "Production OS - Only K300 modules included"
      find $ROOTFS_DIR/lib/modules | grep "\.ko$" | grep  -v -f $BUILD_DIR/config/k300-modules.txt | xargs rm
      ;;
    "MAGELIS")
      echo "Production OS - Only MAGELIS modules included"
      find $ROOTFS_DIR/lib/modules | grep "\.ko$" | grep  -v -f $BUILD_DIR/config/magelis-modules.txt | xargs rm
      ;;
  esac
  
}

build_initrmfs(){
  # Remove the fake initrmfs.cpio before creating the real one from the ROOTFS
  rm -rf $BUILD_DIR/initrmfs.cpio

  # Create the initrmfs from the rootfs directory
  cd $ROOTFS_DIR
  # Cleanup rootfs before creating the initrmfs.cpio
  find . -type f -name '.empty' -size 0c -delete
  rm -rf usr/man usr/share/man usr/local/man usr/local/share/man
  rm -rf usr/lib/pkgconfig usr/local/lib/pkgconfig
  rm -rf usr/include usr/local/include
  
  # remove static libraries and archives left over from the build
  find $ROOTFS_DIR | grep ".\.la$" | xargs rm
  find $ROOTFS_DIR | grep ".\.a$" | xargs rm
  find $ROOTFS_DIR | grep ".\.o$" | xargs rm
  
  find $ROOTFS_DIR -executable -type f | grep -v '.\.script$' | grep -v '.\.sh$' | grep -v '.\.bin$' | grep -v '.\.txt$'| xargs strip

  find . -print0 | cpio --null --create --verbose  --format=newc > $BUILD_DIR/initrmfs.cpio
}

build_kernel() {
  PS4="[build_kernel] "

  cd $SRC_DIR/linux-$KERNEL_VERSION
  cp -f $BUILD_DIR/config/$KERNEL_CONFIG .config

  make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS
  
  cp arch/x86_64/boot/bzImage $OUT_DIR/BOOTx64.EFI
}

###############################################
# Build Custom Init (Golang)
###############################################

build_custom_init() {
  PS4="[build_custom_init] "

  cd $BUILD_DIR/init/init-custom
  go get github.com/u-root/u-root
  go build -ldflags "-s -w" -o $ROOTFS_DIR/init-custom

  strip $ROOTFS_DIR/init-custom
}

###############################################
# Download all pacakges 
###############################################

download_packages() {
  PS4="[download_packages] "
  download_kernel
  download_wg_tools

  download_nftables
  download_musl
  download_busybox
  download_docker

  touch $SRC_DIR/flag_downloaded
}

###############################################
# Build all pacakges 
###############################################

build_packages() {
  PS4="[build_packages] "
  build_busybox
  build_musl
  build_nftables
  build_wg_tools
  
  touch $SRC_DIR/flag_built_packages
}

install_packages() {
  PS4="[install_packages] "
  install_busybox
  install_musl
  install_nftables
  install_wg_tools
  install_docker
}

build_busybox() {
  cd $SRC_DIR/busybox-$BUSYBOX_VERSION
  make distclean defconfig #-j $NUM_JOBS
  cp -f $BUILD_DIR/config/busybox-config .config
  make oldconfig
}

install_busybox() {
  cd $SRC_DIR/busybox-$BUSYBOX_VERSION
  make EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" busybox install #-j $NUM_JOBS
}

# cheat! using the prebuid alpine docker libs and binaries for cryptscript and e2fsprogs
install_deployment_tools() {
  cp -av /sbin/cryptsetup $ROOTFS_DIR/sbin
  cp -av /lib/libcryptsetup.so.* $ROOTFS_DIR/lib
  cp -av /lib/libpopt.so.* $ROOTFS_DIR/lib
  cp -av /lib/libuuid.so.* $ROOTFS_DIR/lib
  cp -av /lib/libblkid.so.* $ROOTFS_DIR/lib 
  cp -av /lib/libdevmapper.so.1.* $ROOTFS_DIR/lib
  cp -av /lib/libcrypto.so.1.* $ROOTFS_DIR/lib 
  cp -av /usr/lib/libargon2.so.* $ROOTFS_DIR/usr/lib
  cp -av /usr/lib/libjson-c.so.* $ROOTFS_DIR/usr/lib

  cp -av /sbin/mke2fs $ROOTFS_DIR/sbin

  cp -av /lib/libext2fs.so.* $ROOTFS_DIR/lib 
  cp -av /lib/libcom_err.so.* $ROOTFS_DIR/lib 
  cp -av /lib/libblkid.so.* $ROOTFS_DIR/lib 
  cp -av /lib/libuuid.so.* $ROOTFS_DIR/lib 
  cp -av /lib/libe2p.so.* $ROOTFS_DIR/lib 

  cd $ROOTFS_DIR
  ln -sfv mke2fs sbin/mkfs.ext4 
  ln -sfv mke2fs sbin/mkfs.ext3
  
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

  if [ ! -d $ROOTFS_DIR ]; then
    mkdir -p $ROOTFS_DIR
  fi
  rm -rf $ROOTFS_DIR/*
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

  if [ ! -f $SRC_DIR/flag_built_modules ]; then
    build_modules
  fi

  # build the Golang init command
  build_custom_init
  # build the rootfs
  build_rootfs

  install_packages
  install_modules
  
case "$BUILD_TYPE" in
  "FACTORY")
    echo "including the Factory EFI Image"
    install_deployment_tools
    ;;
  *)
    echo "Production OS for $BUILD_TYPE"
    rm $ROOTFS_DIR/sbin/*.sh
    ;;
esac

  # strip and create the initramfs for the kernel build 
  build_initrmfs
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
  BUILD_TYPE=$1
  echo "Building EFI for: $BUILD_TYPE"
  build_all
  echo "Production OS build completed for: $BUILD_TYPE"

  ;;
esac