#!/bin/bash

# SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: Apache-2.0

set -ex

# build_rootfs(OUTPUT_DIRECTORY, TEMPLATE_DIRECTORY)
build_rootfs() {

	OUT=${1}
	TEMPLATE=${2}

	mkdir -p ${OUT}
	cp -a ${TEMPLATE}/* ${OUT}

	# template directories contain a file called .empty, to ensure git doesn't ignore the empty directory
	find ${OUT} -type f -name '.empty' -size 0c -delete

	# Set permissions
	cd ${OUT} && chmod 0755 bin dev etc proc sbin sys usr usr/bin usr/sbin
	cd ${OUT} && chmod -R 0777 tmp var
	cd ${OUT} && chmod 0770 root

	cd ${OUT} && chmod 0644 etc/localtime
	cd ${OUT} && chmod 0644 etc/resolv.conf

	# create temp character devices to allow for inital boot
	cd ${OUT} && mknod dev/console c 5 1
	cd ${OUT} && chmod 0600 dev/console
	cd ${OUT} && mknod dev/tty c 5 0
	cd ${OUT} && chmod 0666 dev/tty
	cd ${OUT} && mknod dev/null c 1 3
	cd ${OUT} && chmod 0666 dev/null
	cd ${OUT} && mknod dev/port c 1 4
	cd ${OUT} && chmod 0640 dev/port
	cd ${OUT} && mknod dev/urandom c 1 9
	cd ${OUT} && chmod 0640 dev/urandom

	cd ${OUT} && chmod 770 usr/share/udhcpc/default.script

	# Remove filler
	cd ${OUT} && rm -rf usr/man usr/share/man usr/local/man usr/local/share/man
	cd ${OUT} && rm -rf usr/lib/pkgconfig usr/local/lib/pkgconfig
	cd ${OUT} && rm -rf usr/include usr/local/include
}

build_rootfs $*
