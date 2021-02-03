#!/bin/bash
set -e

ARG1="$1"
ARG2="$2"
ARG3="$3"
ARG4="$4"

MEM="2048"
MNTDIR="/tmp/qemu/mnt"

init() {
    echo ""
    echo "INIT"

    rm -rf $MNTDIR
    mkdir -p $MNTDIR/a
    mkdir -p $MNTDIR/b
    mkdir -p $MNTDIR/c

    TARGET="$ARG2"
    EFI="$ARG3"
    CONFIG="$ARG4"

    if [ ! -f $EFI ]; then
        echo "Error: EFI $EFI does not exist"
        exit
    fi

    if [ "$(uname)" = "Linux" ]; then
        RUNTIME="Linux"
    else
        echo "Error: Unsupported platform: $(uname)"
        exit
    fi

    if [ "$ARG1" = "kvm" ]; then
        QEMU_ACCEL="--enable-kvm -cpu host -smp 4 -machine type=q35,accel=kvm"
        QEMU_DISPLAY="-display gtk -vga std"
    fi

    if [ "$ARG1" = "kvm-headless" ]; then
        QEMU_ACCEL="--enable-kvm -cpu host -smp 4 -machine type=q35,accel=kvm"
        QEMU_DISPLAY="-display curses"
    fi

    if [ "$ARG1" = "emulate" ]; then
        QEMU_ACCEL=""
        QEMU_DISPLAY="-display gtk -vga std"
    fi

    if [ "$ARG1" = "emulate-headless" ]; then
        QEMU_ACCEL=""
        QEMU_DISPLAY="-display curses"
    fi

    QEMU_BIOS="-bios /usr/share/qemu/OVMF.fd"

    echo "Initialized for $RUNTIME with options [$QEMU_DISPLAY $QEMU_ACCEL] using bios: $QEMU_BIOS, baseimage: $TARGET, efi: $EFI, config: $CONFIG"
}

image() {
    echo ""
    echo "IMAGE"

    if [ ! -f $TARGET ]; then
        echo "Baseimage not found, creating new"
        qemu-img create -f raw $TARGET 1000M
    fi

    echo "Partitioning disk"
    echo '' | sudo sfdisk $TARGET <<EOF
    start=2048, type=ef, size=398MiB
    start=400MiB, size=100MiB
    start=500MiB, size=500MiB
EOF

    LOOPDEV=$(sudo losetup --find --show $TARGET)
    sudo partprobe ${LOOPDEV}

    echo "Creating filesystems"
    sudo mkfs.fat -F32 ${LOOPDEV}p1
    sudo mkfs.ext4 -F ${LOOPDEV}p2
    sudo mkfs.ext4 -F ${LOOPDEV}p3

    echo "Inserting EFI"
    sudo mount ${LOOPDEV}p1 $MNTDIR/a
    sudo mkdir -p $MNTDIR/a/EFI/BOOT/
    sudo cp $EFI $MNTDIR/a/EFI/BOOT/BOOTx64.EFI

    echo "Inserting secondary config"
    sudo mount ${LOOPDEV}p2 $MNTDIR/b
    sudo cp $CONFIG $MNTDIR/b/secondary.json
}

verify() {
    echo ""
    echo "VERIFY"

    LOOPDEV=$(sudo losetup --find --show $TARGET)
    sudo partprobe ${LOOPDEV}
    sudo mount ${LOOPDEV}p1 $MNTDIR/a
    sudo mount ${LOOPDEV}p2 $MNTDIR/b
    sudo mount ${LOOPDEV}p3 $MNTDIR/c
    tree $MNTDIR
}

cleanup() {
    echo ""
    echo "CLEANUP"
    sudo umount $MNTDIR/a || true
    sudo umount $MNTDIR/b || true
    sudo umount $MNTDIR/c || true

    LOOPDEV=$(sudo losetup --find --show $TARGET)
    sudo losetup -d ${LOOPDEV} || true
    rm -rf $MNTDIR
    mkdir $MNTDIR
    mkdir $MNTDIR/a
    mkdir $MNTDIR/b
    mkdir $MNTDIR/c
}

run() {
    echo ""
    echo "RUN"
    qemu-system-x86_64 $QEMU_ACCEL $QEMU_DISPLAY $QEMU_BIOS \
        -m $MEM \
        -drive format=raw,file=$TARGET,if=none,id=os2 \
        -device ich9-ahci,id=ahci \
        -device ide-drive,drive=os2,bus=ahci.0 \
        -device virtio-rng-pci \
        -nodefaults \
        -boot c
}

init
cleanup
image
cleanup
verify
cleanup
run
