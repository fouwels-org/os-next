#!/bin/sh
DEVFILE=/tmp/vfat.txt
DEV="/dev/sda1" # default vfat location

if [ -f "$DEVFILE" ]; then
    DEV=`cat /tmp/vfat.txt`
    echo "partition obtained from the tmp filesystem : $DEV"

fi

echo "BOOTx64 is on partition $DEV"

FILE=/mnt
if [ ! -d "$FILE" ]; then
    mkdir /mnt
    echo "$FILE did not exist, but has been created"
fi

echo "downloading EFI image"
wget -O /tmp/BOOTx64.EFI http://81.201.135.86:11223/v1/files/Qk9PVHg2NC5FRkk=
echo "download complete"
if [ $? -eq 0 ]; then
    mount -t vfat $DEV /mnt
    rm /mnt/EFI/BOOT/BOOTx64.EFI
    cp /tmp/BOOTx64.EFI /mnt/EFI/BOOT/BOOTx64.EFI
    umount /mnt
    rm /tmp/BOOTx64.EFI
    echo "BOOTx64.EFI has been upgraded"
else
    echo "Error downloading file"
fi