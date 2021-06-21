#!/bin/bash
set -ex

COMMAND="$1"
MODE="$2"
TARGET="$3"
EFI="$4"

TPM_PATH=/tmp/vtpm

# start starts the software TPM device
# requires https://github.com/stefanberger/swtpm to be installed
startTPM() {

    mkdir -p $TPM_PATH
    echo "** creating vTPM in $TPMDIR **"
    swtpm socket --tpm2 --tpmstate dir=$TPM_PATH --ctrl type=unixio,path=$TPM_PATH/socket --log level=20
    echo "** vTPM closed **"
    echo ""
}

deploy() {   

    if [ "$MODE" = "clean" ]; then
        rm -rf $TARGET || true
    fi

    if [ ! -f $TARGET ]; then
        echo "Baseimage not found, creating new"
        qemu-img create -f raw $TARGET 5000M

        # Initialise disk as GPT
        echo '' | sudo sfdisk $TARGET <<EOF
            label: gpt
            label-id: C8C6B4E3-B2B8-574A-A78B-713A9E8D6013
            start=2048, size=405MiB, type=C12A7328-F81F-11D2-BA4B-00A0C93EC93B, name=60D0-7209
            size=500MiB, type=0FC63DAF-8483-4772-8E79-3D69D8477DE4, name=8cda03ee-602d-4ad6-996b-93c952d88b54
            size=1000MiB, type=0FC63DAF-8483-4772-8E79-3D69D8477DE4, name=8de49650-c95e-433b-ac23-99aeaa726077
EOF
        LOOPDEV=$(sudo losetup --find --show $TARGET)
        sudo partprobe ${LOOPDEV}
        sudo mkdosfs -n "BOOT" ${LOOPDEV}p1
        sudo mke2fs -t "ext4" -L "CONFIG" ${LOOPDEV}p2
        sudo mke2fs -t "ext4" -L "DATA" ${LOOPDEV}p3
        sudo losetup -d ${LOOPDEV}
    fi

    if [ ! -f $EFI ]; then
        echo "Error: EFI $EFI does not exist"
        exit
    fi

    if [ "$COMMAND" = "disk" ]; then
        # insert EFI
        LOOPDEV=$(sudo losetup --find --show $TARGET)
        sudo partprobe ${LOOPDEV}
        mkdir -p /tmp/a
        sudo mount ${LOOPDEV}p1 /tmp/a
        sudo mkdir -p /tmp/a/EFI/BOOT/
        sudo cp $EFI /tmp/a/EFI/BOOT/BOOTx64.EFI
        sudo umount /tmp/a
        rm -rf /tmp/a
        sudo losetup -d ${LOOPDEV}
    fi

    if [ "$(uname)" = "Linux" ]; then
        RUNTIME="Linux"
    else
        echo "Error: Unsupported platform: $(uname)"
        exit
    fi

    QEMU_DISPLAY="-display vnc=:0"

    if [ "$COMMAND" = "kernel" ]; then
        QEMU_DISK="-kernel $EFI --append console=ttyS0 -nographic"
    fi
    if [ "$COMMAND" = "disk" ]; then
        QEMU_DISK="-serial mon:stdio -boot menu=on\
        -drive if=pflash,format=raw,readonly=on,file=OVMF_CODE.fd 
        -drive if=pflash,format=raw,file=OVMF_VARS.fd"
    fi

    echo "starting QEMU"
    qemu-system-x86_64 $QEMU_DISPLAY $QEMU_DISK \
        --enable-kvm \
        -machine type=q35,accel=kvm \
        -cpu host \
        -smp 4 \
        -m 3072 \
        -vga std \
        -drive file=$TARGET,format=raw,format=raw,if=none,id=os2 \
        -device ich9-ahci,id=ahci \
        -device virtio-rng-pci \
        -device e1000e,netdev=n1 \
        -device nvme,drive=os2,serial=nvme-1 \
        -netdev user,id=n1 
#       -device tpm-tis,tpmdev=tpm0 \
#       -chardev socket,id=chrtpm,path=$TPM_PATH/socket \
#       -tpmdev emulator,id=tpm0,chardev=chrtpm

}

case $COMMAND in
kernel)
    deploy
    ;;
disk)
    deploy
    ;;
startTPM)
    while true; do
        startTPM
    done
    ;;
esac
