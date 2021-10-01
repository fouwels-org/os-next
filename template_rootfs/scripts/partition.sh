#!/bin/sh
case "$1" in
"K300")
  DEV="/dev/sda" # default vfat location
  BOOT="${DEV}1"
  CONFIG="${DEV}2"
  DATA="${DEV}3"
  ;;
"K700")
  DEV="/dev/nvmen0" # default vfat location
  BOOT="${DEV}p1"
  CONFIG="${DEV}p2"
  DATA="${DEV}p3"
  ;;
"MAGELIS")
  DEV="/dev/sda" # default vfat location
  BOOT="${DEV}1"
  CONFIG="${DEV}2"
  DATA="${DEV}3"
  ;;
*)
  echo "Please parse the device name you wish to configure"
  echo "Onlogix Karbon 300 use : K300 "
  echo "Onlogix Karbon 700 use : K700 "
  echo "Schneider Electric Magelis use : MAGELIS"
  exit 1
  ;;
esac

echo "Setting up dhcp on ETH0"
/sbin/ip link set dev lo up
/sbin/ip link set dev eth0 up
/sbin/udhcpc -b -i eth0 -p /var/run/udhcpc.pid

echo "Partition device is: $DEV"

echo "BOOT: $BOOT"
echo "CONFIG: $CONFIG"
echo "DATA: $DATA"

(
  echo o     # clear all existing partitions
  echo n     # create a new partition
  echo p     # primary
  echo 1     # partition number
  echo       # default starting seqment
  echo +512M # size in MB
  echo t     # change the type
  echo ef    # VFAT EFI type
  echo n     # create a new partition
  echo p     # primary
  echo 2     # partition number
  echo       # default starting seqment
  echo +512M # size in MB
  echo n     # create a new partition
  echo p     # primary
  echo 3     # partition number
  echo       # default starting seqment
  echo       # default ending seqment (remainder)
  echo w     # write the partition table
) | fdisk $DEV

echo "Formating $BOOT Partition"
echo y | mkfs.vfat -n BOOT $BOOT
echo "Formating $CONFIG Partition"
echo y | mkfs.ext4 -L CONFIG $CONFIG
echo "Formating $DATA Partition"
echo y | mkfs.ext4 -L DATA $DATA

echo "---------------------"
echo "Setting up os-next"
echo "---------------------"
FILE=/mnt
if [ ! -d "$FILE" ]; then
  mkdir /mnt
  echo "$FILE did not exist, but has been created"
fi

echo "downloading EFI image"
wget -O /tmp/BOOTx64.EFI http://81.201.135.86:11223/v1/files/Qk9PVHg2NC5FRkk=
echo "download complete"
if [ $? -eq 0 ]; then
  mount -t vfat $BOOT /mnt
  mkdir -p /mnt/EFI/BOOT
  cp /tmp/BOOTx64.EFI /mnt/EFI/BOOT/BOOTx64.EFI
  umount /mnt
  rm /tmp/BOOTx64.EFI
  echo "os-next rt OS has been successfully installed"
else
  echo "Error downloading file"
fi
