#!/bin/bash
set -e

ARG1="$1"
ARG2="$2"
ARG3="$3"

MEM="2048"
MNTDIR="/tmp/qemu/mnt"
LOOP="loop9"

init() {
    echo ""
    echo "INIT"

    mkdir -p $MNTDIR

    TARGET="$ARG2"
    EFI="$ARG3"

    if [ ! -f $TARGET ]; then
        echo "Error: Base image $TARGET does not exist"
        exit
    fi

    if [ ! -f $EFI ]; then
        echo "Error: EFI $EFI does not exist"
        exit
    fi

    if [ "$ARG1" = "headless" ]; then
        QEMU_DISPLAY="-display curses"
    else
        QEMU_DISPLAY="-display gtk -vga cirrus"
    fi

    if [ "$ARG1" = "redisk" ]; then
        redisk
        exit
    fi
    if [ "$(uname)" == "Linux" ]; then
        RUNTIME="Linux"
        QEMU_KVM="--enable-kvm -cpu host -smp 4 -machine type=q35,accel=kvm"
        QEMU_BIOS="/usr/share/qemu/OVMF.fd"
    else
        echo "Error: Unsupported platform: $(uname)"
        exit
    fi

    echo "Initialized for $RUNTIME with options [$QEMU_DISPLAY $QEMU_KVM] using bios: $QEMU_BIOS, baseimge: $TARGET, efi: $EFI"
}

redisk() {
    echo ""
    echo "REDISK"
    qemu-img create -f raw $TARGET 1124M
}

image() {
    echo ""
    echo "IMAGE"

    sudo kpartx -s -a -v $TARGET
    sudo mount /dev/mapper/${LOOP}p1 /$MNTDIR

    sudo mkdir -p $MNTDIR/EFI/boot/
    sudo cp $EFI $MNTDIR/EFI/boot/bootx64.efi

    sudo umount $MNTDIR
    sudo kpartx -s -d -v $TARGET

    # Verify
    sudo kpartx -s -a -v $TARGET
    sudo mount /dev/mapper/${LOOP}p1 $MNTDIR

    tree $MNTDIR
}

cleanup() {
    echo ""
    echo "CLEANUP"
    sudo umount $MNTDIR || true
    sudo kpartx -s -d -v $TARGET || true
}

run() {
    echo ""
    echo "RUN"
    qemu-system-x86_64 $QEMU_KVM $QEMU_DISPLAY \
        -m $MEM \
        -drive format=raw,file=$TARGET \
        -bios $QEMU_BIOS \
        -nodefaults \
        -boot c
}

init
image || true
cleanup
run
