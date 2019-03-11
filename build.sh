#!/bin/sh

set -ex

KERNEL_VERSION=4.20.12
MUSL_VERSION=1.1.21
BUSYBOX_VERSION=1.30.1
DROPBEAR_VERSION=2018.76
RNGTOOLS_VERSION=5
SYSLINUX_VERSION=6.03
IPTABLES_VERSION=1.8.2
DOCKER_VERSION=18.09.2
KMOD=26

NUM_JOBS="$(grep ^processor /proc/cpuinfo | wc -l)"

build=/build
rootfs=$build/rootfs/
isoimage=$build/isoimage/

debug() { echo "Dropping into a shell for debugging ..."; /bin/sh; }

config() { 
  if grep "CONFIG_$2" .config; then
    sed -i "s|.*CONFIG_$2.*|CONFIG_$2=$1|" .config
  else
    echo "CONFIG_$2=$1" >> .config
  fi
}

download_kmod() {
  cd /build/src
  wget -q -O kmod.tar.xz \
    https://mirrors.edge.kernel.org/pub/linux/utils/kernel/kmod/kmod-$KMOD.tar.xz
  tar -xf kmod.tar.xz
}

download_syslinux() {
  cd /build/src
  wget -q -O syslinux.tar.xz \
    http://kernel.org/pub/linux/utils/boot/syslinux/syslinux-$SYSLINUX_VERSION.tar.xz
  tar -xf syslinux.tar.xz
}

download_systemd_boot() {
  cd /build/src
  wget -q -O systemd_boot.tar.xz \
    https://github.com/ivandavidov/systemd-boot/releases/download/systemd-boot_26-May-2018/systemd-boot_26-May-2018.tar.xz
  tar -xf systemd_boot.tar.xz
}

download_kernel() {
  cd /build/src
  wget -q -O kernel.tar.xz \
    http://kernel.org/pub/linux/kernel/v4.x/linux-$KERNEL_VERSION.tar.xz
  tar -xf kernel.tar.xz
}

download_musl() {
  cd /build/src    
  wget -q -O musl.tar.gz \
    http://www.musl-libc.org/releases/musl-$MUSL_VERSION.tar.gz
  tar -xf musl.tar.gz
}

download_busybox() {
  cd /build/src    
  wget -q -O busybox.tar.bz2 \
    http://busybox.net/downloads/busybox-$BUSYBOX_VERSION.tar.bz2
  tar -xf busybox.tar.bz2
}

download_rngtools() {
  cd /build/src     
  #wget -q -O rngtools.tar.gz https://sourceforge.net/projects/gkernel/files/rng-tools/5/rng-tools-5.tar.gz
  cp ../source/rng-tools-5.tar.gz rngtools.tar.gz
  tar -xf rngtools.tar.gz
}

download_dropbear() {
  cd /build/src    
  wget -q -O dropbear.tar.bz2 \
    https://matt.ucc.asn.au/dropbear/dropbear-$DROPBEAR_VERSION.tar.bz2
  tar -xf dropbear.tar.bz2
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
  mkdir -p proc
  mkdir -p run
  mkdir -p sys
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
  cp cloudinit /build/rootfs/usr/bin
  cp default.script /build/rootfs/usr/share/udhcpc
  cp lastlog /build/rootfs/var/log
  cp utmp /build/rootfs/var/run
  cp mail /build/rootfs/var
  cp ca-bundle.pem cert.pem /build/rootfs/etc/ssl/certs
  cp init /build/rootfs/
  )
}

build_musl() {
  (
  cd musl-$MUSL_VERSION
  ./configure \
    --prefix=/usr
  make
  make DESTDIR=$rootfs install
  )
}

  
build_busybox() {
  (
  cd busybox-$BUSYBOX_VERSION
  make distclean defconfig #-j $NUM_JOBS
  config y STATIC
  config n INCLUDE_SUSv2
  config y INSTALL_NO_USR
  config "\"$rootfs\"" PREFIX
  config y FEATURE_EDITING_VI
  config y TUNE2FS
  config n BOOTCHARTD
  config n INIT
  config n LINUXRC
  config y FEATURE_GPT_LABEL
  config n LPD
  config n LPR
  config n LPQ
  config n RUNSV
  config n RUNSVDIR
  config n SV
  config n SVC
  config n SVLOGD
  config n HUSH
  config n CHAT
  config n CONSPY
  config n RUNLEVEL
  config n PIPE_PROGRESS
  config n RUN_PARTS
  config n START_STOP_DAEMON
  yes "" | make oldconfig
  make \
    EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" \
    busybox install #-j $NUM_JOBS
  )
}

build_rngtools() {
  (
  cd rng-tools-$RNGTOOLS_VERSION
  ./configure \
    --prefix=/usr \
    --sbindir=/usr/sbin \
    CFLAGS="-static" LIBS="-l argp"
  make
  make DESTDIR=$rootfs install
  )
}

build_dropbear() {
  (
  cd dropbear-$DROPBEAR_VERSION
  ./configure \
    --prefix=/usr \
    --mandir=/usr/man \
    --enable-static \
    --disable-zlib \
    --disable-wtmp \
    --disable-syslog

  make \
    EXTRA_CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" \
    DESTDIR=$rootfs \
    PROGRAMS="dropbear dbclient dropbearkey scp" \
    strip install #-j $NUM_JOBS
  ln -sf /usr/bin/dbclient $rootfs/usr/bin/ssh
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
    #-j $NUM_JOBS
  make DESTDIR=$rootfs install
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

install_docker() {
  mv docker/* $rootfs/usr/bin/
}

write_metadata() {
  cat > $rootfs/etc/os-release <<EOF
NAME=EniCore
VERSION=V0.0.1
ID=enicore
ID_LIKE=tcl
VERSION_ID=V0.0.1
PRETTY_NAME="EniCore Linux V0.0.1 (TCL 0); "
ANSI_COLOR="1;34"
HOME_URL="http://lagoni.co.uk/"
SUPPORT_URL="http://lagoni.co.uk"
BUG_REPORT_URL="http://lagoni.co.uk"
EOF

  cat > $rootfs/usr/bin/mcl <<EOF
#!/bin/sh

echo EniCore Linux (ECL) V0.0.1"

# End of file
EOF
chmod +x $rootfs/usr/bin/mcl
}

build_rootfs() {
  (
  cd $rootfs

  # Cleanup rootfs
  find . -type f -name '.empty' -size 0c -delete
  rm -rf usr/man usr/share/man
  rm -rf usr/lib/pkgconfig
  rm -rf usr/include

  # Archive rootfs
  find . | cpio -R root:root -H newc -o | gzip -9 > /build/src/rootfs.gz
  )
}

sync_rootfs() {
  (
  mkdir rootfs.old
  cd rootfs.old
  zcat $build/rootfs.gz | cpio -idm
  rsync -aru . $rootfs
  )
}

build_kernel() {
  (
  cd /build/src
  cd linux-$KERNEL_VERSION
  cp /build/config .config
  #make mrproper defconfig -j $NUM_JOBS
  # NOT NEEDED WITH IF THE KERNEL CONFIG IS CORRECRTLY CONFIGURED
  #  yes "" | make oldconfig
  
  make CFLAGS="-Os -s -fno-stack-protector -U_FORTIFY_SOURCE" -j $NUM_JOBS
  make INSTALL_MOD_PATH=$rootfs modules_install
  cp arch/x86_64/boot/bzImage /build/src/kernel.gz
  )
}

check_root() {
  if [ ! "$(id -u)" = "0" ] ; then
    cat << CEOF

  ISO image preparation process for UEFI systems requires root permissions
  but you don't have such permissions. Restart this script with root
  permissions in order to generate UEFI compatible ISO structure.

CEOF
    exit 1
  fi
}

build_boot_uefi() {
  check_root

  test -d "$isoimage" || mkdir "$isoimage"

  MLL_CONF=x86_64

  LOADER=/build/src/systemd-boot*/uefi_root/EFI/BOOT/BOOTx64.EFI

  # Find the kernel size in bytes.
  kernel_size=`stat -c %s /build/src/kernel.gz`

  # Find the initramfs size in bytes.
  rootfs_size=`stat -c %s /build/src/rootfs.gz`

  loader_size=`stat -c %s $LOADER`

  # The EFI boot image is 64KB bigger than the kernel size ( + 2MB to ensure the image is big enough, not sure why it runs out of space otherwise).
  image_size=$((kernel_size + rootfs_size + loader_size + 2097152 + 65536))

  echo "Creating UEFI boot image file 'uefi.img'."
  rm -f uefi.img
  truncate -s $image_size uefi.img

  echo "Attaching hard disk image file to loop device."
  LOOP_DEVICE_HDD=$(losetup -f)
  losetup $LOOP_DEVICE_HDD uefi.img

  echo "Formatting hard disk image with FAT filesystem."
  mkfs.vfat $LOOP_DEVICE_HDD

  echo "Preparing 'uefi' work area."
  rm -rf uefi
  mkdir -p /uefi
  mount uefi.img /uefi

#  # Add the configuration files for UEFI boot.
#  cp -r $SRC_DIR/minimal_boot/uefi/* \
#    $ISOIMAGE

  echo "Preparing kernel and rootfs."
  mkdir -p /uefi/minimal/$MLL_CONF
  cp /build/src/kernel.gz \
    /uefi/minimal/$MLL_CONF/kernel.xz
  cp /build/src/rootfs.gz \
    /uefi/minimal/$MLL_CONF/rootfs.xz

  echo "Preparing 'systemd-boot' UEFI boot loader."
  mkdir -p /uefi/EFI/BOOT
  cp $LOADER \
    /uefi/EFI/BOOT

  echo "Preparing 'systemd-boot' configuration."
  mkdir -p /uefi/loader/entries
  cp /build/src/systemd-boot*/uefi_root/loader/loader.conf \
    /uefi/loader
  cp /build/src/systemd-boot*/uefi_root/loader/entries/mll-${MLL_CONF}.conf \
    /uefi/loader/entries

  echo "Setting the default UEFI boot entry."
  sed -i "s|default.*|default mll-$MLL_CONF|" /uefi/loader/loader.conf

  echo "Unmounting UEFI boot image file."
  sync
  umount /uefi
  sync
  sleep 1 

  # The directory is now empty (mount point for loop device).
  #rm -rf /uefi

  # Make sure the UEFI boot image is readable.
  chmod ugo+r uefi.img

  mkdir -p /build/isoimage/boot
  cp uefi.img \
    /build/isoimage/boot/uefi.img
  cd $isoimage

  # Now we generate 'hybrid' ISO image file which can also be used on
  # USB flash drive, e.g. 'dd if=minimal_linux_live.iso of=/dev/sdb'.
  xorriso -as mkisofs \
    -isohybrid-mbr /build/src/syslinux-*/bios/mbr/isohdpfx.bin \
    -c boot/boot.cat \
    -e boot/uefi.img \
      -no-emul-boot \
      -isohybrid-gpt-basdat \
    -o /build/isoimage/enicore_uefi_0_0.iso \
    $isoimage

}

download_packages(){
  cd /build/src
  download_busybox
  download_kmod
  download_musl
  download_rngtools
  download_dropbear
  download_iptables
  download_kernel

  download_syslinux
  download_systemd_boot

  download_docker
  install_docker
}

build_packages(){
  cd /build/src
  build_busybox
  build_kmod
  build_musl
  build_iptables
  build_rngtools
  build_dropbear
  
}

prepare_rootfs(){
  write_metadata
  build_rootfs
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
  # makes the kernel into  kernel.gz file
  build_kernel  
 
  # WORDIR build/src
  build_packages
  
  # WORDIR $rootfs
  # Creates the release file in the rootfs
  # makes the rootfs into a rootfs.gz file
  prepare_rootfs

  # WORDIR $isoimage  
  build_boot_uefi
}

repack() {
  sync_rootfs
  write_metadata
  build_rootfs
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
  build_uefi)
    prepare_rootfs
    build_boot_uefi
    ;;
  repack)
    repack
    ;;
  *)
    build_all
    ;;
esac
