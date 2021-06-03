# IIoT OS Mjolnir Build (Thor's Hammer)
This is a real-time IIoT OS made for Docker. All is build from within a container using musl rather than libc. The Dockerfile included is used to create the toolchain image. 

It is a EFI based boot system.

## Quick Start

In the `build.sh` script use the KERNEL_CONFIG variable to set the kernel config to be used in the compilation 

Ensure docker is compiling with buildkit for a layer based build log (highly recommended)

    export DOCKER_BUILDKIT=1 

Build the kernel EFI (`make <target> # make fast|k300|schneider|...`)

    docker build -t <build arguments> registry2.lagoni.co.uk/os_build_env:local .

The following build arguments are specified.

_Included for documentatin, see Makefile for existing targets, makefile should be used instead of direct calling_

    --build-arg CONFIG_COMPRESSION=XZ # Kernel/Initramfs compression scheme. (LZ4 for fast, XZ for small)
    --build-arg CONFIG_PRIMARY=nvme.json # One of config/primary
    --build-arg CONFIG_MODULES=ALL # ALL or One of config/modules

Copy the kernel EFI from the built image to ./out (`make run`)

    docker run -it --rm -v $(PWD)/out:/out registry2.lagoni.co.uk/os_build_env:local

The dockerfile is constructed to cache layers between builds, if required files have not been modified. There is no need to preserve files on a volume.

## Deploying the EFI image

The image is copied out of the container when the build is successfully completed and place in a subdirectory called `out` in current working director. The file is called BOOTx64.EFI

To deploy this onto a physical hardware device, this hardware needs to support UEFI boot, which is a modern boot loader, which is suported in most new BIOS implementaitons. 

To make this work, first format a drive (can be USB or HDD) with FAT32. Then create a directory structure on this device /EFI/BOOT and simply copy the BOOTx64.EFI into the folder. The BIOS will identify this path and boot.

See `deploy/qemu` for a software deployment. This will automatically format and set up a drive, before starting with QEMU/KVM.

Unless specified, run `make qemu-kernel` within, to build and start the kernel and connect stdin/stdout over virtual IO.

## CHECKS

`registry2.lagoni.co.uk/kernel-check:latest`

## TIPS and TRICKS

In the container you can find the linux source directory under /build/src/linux... 

Use `menuconfig` to setup the kernel then copy the .config to /build/out 

    `cp .config /build/out`

This will copy the .config to the host machine.

## License
APACHE 2.0
