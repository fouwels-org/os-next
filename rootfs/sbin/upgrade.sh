#!/bin/sh
mount -t vfat /dev/sda1 /mnt

rm /mnt/EFI/BOOT/BOOTx64.EFI

#wget -O bzImage http://81.201.135.86:11223/v1/files/Qk9PVHg2NC5FRkk=
#wget -O initrmfs.cpio.xz http://81.201.135.86:11223/v1/files/aW5pdHJtZnMuY3Bpby54eg==

wget -O BOOTx64.EFI http://81.201.135.86:11223/v1/files/Qk9PVHg2NC5FRkk=

#wget -O initramfs.cpio http://81.201.135.86:11223/v1/files/aW5pdHJtZnMuY3Bpbw==


umount /mnt
