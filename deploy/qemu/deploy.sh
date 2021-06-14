#!/bin/bash
set -e

MODE="$1"
TARGET="$2"
EFI="$3"

run() {
    echo "RUN"

    QEMU_DISPLAY="-display vnc=:0"
    QEMU_DISK="-kernel $EFI --append console=ttyS0 -nographic"

    if [ "$MODE" = "clean" ]; then
        rm -rf $TARGET || true
    fi

    if [ ! -f $TARGET ]; then
        echo "Baseimage not found, creating new"
        qemu-img create -f raw $TARGET 1500M
    fi

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
        -netdev user,id=n1

}
run
