#!/bin/sh
mkdir -p /root/usb
mkdir -p /root/mnt
mount -t vfat /dev/sda1 /root/mnt
mount -t vfat /dev/sdb1 /root/usb
FILE="/root/usb/EFI/BOOT/BOOTx64.EFI"
if [ -f $FILE ]; then
    echo "Deleting existing file from USB stick: $FILE"
    rm -f $FILE
fi

cp /root/mnt/EFI/BOOT/BOOTx64.EFI /root/usb/EFI/BOOT
umount /root/mnt
umount /root/usb
