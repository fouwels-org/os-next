<!--
SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>

SPDX-License-Identifier: Apache-2.0
-->

# os-next
This is a real-time IIoT OS made for Docker. All is build from within a container using musl rather than libc. The Dockerfile included is used to create the toolchain image. 

It is a EFI based boot system.

## Quick Start

In the `build.sh` script use the KERNEL_CONFIG variable to set the kernel config to be used in the compilation 

Ensure docker is compiling with buildkit for a layer based build log (highly recommended)

    export DOCKER_BUILDKIT=1 

Build the kernel EFI (`make <target> # make k300|schneider|...`)

    docker build -t <build arguments> os-next

The following build arguments are specified.

_Included for documentation, see Makefile for existing targets, makefile should be used instead of direct calling_

    --build-arg CONFIG_PRIMARY=standard.yml # One of config/primary
    --build-arg CONFIG_MODULES=ALL # ALL or One of config/modules
    --build-arg COMPRESSION_LEVEL=9 # (optionally) override the default kernel ZSTD compression level (9 for fast, 22 for maximum)

Copy the kernel EFI from the built image to ./out (`make run`)

    docker run -it --rm -v $(PWD)/out:/out os-next

The dockerfile is constructed to cache layers between builds, if required files have not been modified. There is no need to preserve files on a volume.

## Deploying the EFI image

The image is copied out of the container when the build is successfully completed and place in a subdirectory called `out` in current working director. The file is called BOOTx64.EFI

To deploy this onto a physical hardware device, this hardware needs to support UEFI boot, which is a modern boot loader, which is suported in most new BIOS implementaitons. 

To install, format a drive with the labels and filesystems stated in the primary config.

Then create a directory structure on this device <partition 1>/EFI/BOOT and simply copy the BOOTx64.EFI into the folder. The UEFI will identify this path and boot.

The OS will use the LABEL fields to map devices to mount points specified in the primary config. Device IDs (eg. /dev/sda1) are not used to allow common operation across disk types (eg. /dev/nvme0n1) 

See `deploy/qemu` for a software deployment. This will automatically format and set up a drive, before starting with QEMU/KVM.

Unless specified, run `make deploy-clean` within, to build and start the kernel and connect stdin/stdout over virtual IO.

## TIPS and TRICKS

In the container you can find the linux source directory under /build/src/linux... 

Use `menuconfig` to setup the kernel then copy the .config to /build/out 

    `cp .config /build/out`

This will copy the .config to the host machine.
