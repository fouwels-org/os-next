#!/bin/bash
set -e

ARG1="$1"
ARG2="$2"
ARG3="$3"
ARG4="$4"

TARGET="$ARG2"
EFI="$ARG3"
CONFIG="$ARG4"

MNTDIR="/tmp/qemu/mnt"

init() {
    echo ""
    echo "INIT"

    rm -rf $MNTDIR
    mkdir -p $MNTDIR/a
    mkdir -p $MNTDIR/b
    mkdir -p $MNTDIR/c

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

    LOOPDEV=$(sudo losetup --find --show $TARGET) || true
    sudo losetup -d ${LOOPDEV} || true
    rm -rf $MNTDIR
    mkdir -p $MNTDIR
    mkdir -p $MNTDIR/a
    mkdir -p $MNTDIR/b
    mkdir -p $MNTDIR/c
}

run() {
    echo ""
    echo "RUN"

    if [ "$ARG1" = "kvm-gtk" ]; then
        QEMU_DISPLAY="-display gtk"
        QEMU_DISK="-nodefaults -boot c -bios /usr/share/ovmf/bios.bin -monitor stdio"
    fi

    if [ "$ARG1" = "kvm-vnc" ]; then
        QEMU_DISPLAY="-display vnc=:0"
        QEMU_DISK="-nodefaults -boot c -bios /usr/share/ovmf/bios.bin -monitor stdio"
    fi

    if [ "$ARG1" = "kvm-kernel" ]; then
        QEMU_DISPLAY="-display vnc=:0"
        QEMU_DISK="-kernel $EFI --append console=ttyS0 -nographic"
    fi
    
    qemu-system-x86_64 $QEMU_DISPLAY $QEMU_DISK \
        --enable-kvm \
        -machine type=q35,accel=kvm \
        -cpu host \
        -smp 4 \
        -m 3072 \
        -vga std \
        -drive format=raw,file=$TARGET,if=none,id=os2 \
        -device ich9-ahci,id=ahci \
        -device nvme,drive=os2,serial=nvme-1 \
        -device virtio-rng-pci \
        -device e1000e,netdev=n1 \
        -netdev user,id=n1 \
        
}

init
cleanup
image
cleanup
verify
cleanup
run
