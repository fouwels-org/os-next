#!/bin/bash
set -ex

# build_deptools(OUTPUT_DIRECTORY, ENABLED)
build_deptools() {

	OUT=${1}
	ENABLED=${2}

	if [ ${ENABLED} -eq "1" ]; then

		cp -av /sbin/cryptsetup ${OUT}/sbin
		cp -av /lib/libcryptsetup.so.* ${OUT}/lib
		cp -av /lib/libpopt.so.* ${OUT}/lib
		cp -av /lib/libuuid.so.* ${OUT}/lib
		cp -av /lib/libblkid.so.* ${OUT}/lib
		cp -av /lib/libdevmapper.so.1.* ${OUT}/lib
		cp -av /lib/libcrypto.so.1.* ${OUT}/lib
		cp -av /usr/lib/libargon2.so.* ${OUT}/usr/lib
		cp -av /usr/lib/libjson-c.so.* ${OUT}/usr/lib

		cp -av /sbin/mke2fs ${OUT}/sbin

		cp -av /lib/libext2fs.so.* ${OUT}/lib
		cp -av /lib/libcom_err.so.* ${OUT}/lib
		cp -av /lib/libblkid.so.* ${OUT}/lib
		cp -av /lib/libuuid.so.* ${OUT}/lib
		cp -av /lib/libe2p.so.* ${OUT}/lib

		ln -sfv ${OUT}/mke2fs ${OUT}/sbin/mkfs.ext4
		ln -sfv ${OUT}/mke2fs ${OUT}/sbin/mkfs.ext3
	fi
}

build_deptools $*