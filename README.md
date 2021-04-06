# IIoT OS Mjolnir Build (Thor's Hammer)
This is a real-time IIoT OS made for Docker. All is build from within a container using musl rather than libc. The Dockerfile included is used to create the toolchain image. 

It is a EFI based boot system.

## Quick Start

In the `build.sh` script use the KERNEL_CONFIG variable to set the kernel config to be used in the compilation 

Clone the directory into a folder on your local machine

    docker build -t registry.lagoni.co.uk/os_build_env:latest .

Run the toolchain on a host machine. This will give you access to the linux command line in the Alpine toolchain image. 

    docker run -it --rm -v build_data:/build -v $PWD/out:/build/out --privileged=true --name toolchain registry.lagoni.co.uk/os_build_env:latest /bin/ash

Compile go init from the host volume. 

    docker run -it --rm -v build_data:/build -v $PWD/out:/build/out -v $PWD/init:/build/init -v $PWD/config:/build/config -v $PWD/scripts:/build/scripts --name toolchain registry.lagoni.co.uk/os_build_env:latest /bin/ash

At this point you will be taken to the command line within the container. The build process will use the config-enicore as the config file for the .config of the kernel. If you want to configure your own kernel, then change the config-enicore. The kernel version supported is 4.20. 
  
To start the build process run:
  
    /build.sh  - will build all
    /build.sh build_init  - will only build the go init and uinit programs
    /build.sh rebuild  - will only build the go init and uinit programs and the reassmble the EFI. It will not rebuild the kernel, but /build.sh must have been run for this to work

    ./build.sh <CONFIG> <MODULE>|ALL|FACTORY - where <CONFIG> is the config required, one of config/primary/, and <MODULES> is the kernel module list to supply (not load), ALL, or FACTORY: denoting ALL + extra tooling.

When the build process completes there will be a EFI image, called BOOTx64, in the out folder. 
        
## To clean up

When the compilation has been completed the all the data is stored in a docker volume called build_data. In order to clean all the data, source and build data produced during the ISO build, the simplest way is to delete the docker volume after having exited the toolchain container. 

To exit the toolchain container use the command from within the command line of the running container:

     exit

Run the following command to remove the data.

    docker volume rm build_data

## Deploying the EFI image

The image is copied out of the container when the build is successfully completed and place in a subdirectory called `out` in current working director. The file is called BOOTx64.EFI

To deploy this onto a physical hardware device, this hardware needs to support UEFI boot, which is a modern boot loader, which is suported in most new BIOS implementaitons. 

To make this work, first format a drive (can be USB or HDD) with FAT32. Then create a directory structure on this device /EFI/BOOT and simply copy the BOOTx64.EFI into the folder. The BIOS will identify this path and boot.

See `deploy/qemu` for a software deployment. This will automatically format and set up a drive, before starting with QEMU/KVM.

## CHECKS

`registry2.lagoni.co.uk/kernel-check:latest`

## TIPS and TRICKS

In the container you can find the linux source directory under /build/src/linux... 

Use `menuconfig` to setup the kernel then copy the .config to /build/out 

    `cp .config /build/out`

This will copy the .config to the host machine.

## License
APACHE 2.0
