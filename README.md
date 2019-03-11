# OS Mjolnir Build (Thor's Hammer)

OSNano is a minimalist linux operating system designed to run Docker containers. The operating system build a linux ISO from source. It will download the source dependencies and build these. 

All is build from within a container using musl rather than libc. The Dockerfile included is used to create the toolchain image. 

It is a UEFI based boot system.

## Quick Start

Clone the directory into a folder on your local machine

    git clone https://gitlab.com/lagoni-dev/os-nano.git
    cd os-nano
 Build the Toolchain 

    docker build -t registry.lagoni.co.uk/os_build_env:latest .

Run the toolchain on a host machine. This will give you access to the linux command line in the Alpine toolchain image. 

    docker run -it --rm -v build_data:/build --privileged=true --name toolchain registry.lagoni.co.uk/os_build_env:latest /bin/ash

At this point you will be taken to the command line within the container. The build process will use the config-enicore as the config file for the .config of the kernel. If you want to configure your own kernel, then change the config-enicore. The kernel version supported is 4.20. 
  
To start the build process run:
  
    ./build.sh

When the build process completes there will be a UEFI ISO image, called enicore_uefi_0_0.iso, in the isoimage folder. 

To confirm that the image has been created run the following command:

    ls -lah /build/isoimage/enicore_uefi_0_0.iso 

This command should show an iso image called 

    enicore_uefi_0_0.iso (approx 60Mb in size)

## Getting the ISO image

To copy the UEFI ISO from the toolchain container to the host machine you first need to keep open a second terminal while the toolchain container remains running. 

At the command line back at the host machine (2nd terminal window) run:

    docker cp toolchain:/build/isoimage/enicore_uefi_0_0.iso ~/enicore_uefi_0_0.iso

This uses the the docker client to copy from the runnning toolchain container volume. The command syntax is:

    docker cp <containerId>:/file/path/within/container /host/path/target
        
Alternatively, is to start another container, called the filebrowser, which provides as web interface to the docker volume where the ISO file has been build. To start the filebrowser use the following in a docker-compose file and run

    docker-compose up -d:

Point your webbrowser to localhost:8888 and authenticate with the credentials username: admin, password: admin. This approach binds the same docker volume build_data and exposes the files and directories via the web interface to allow a web browser to download the iso image.

## To clean up

When the compilation has been completed the all the data is stored in a docker volume called build_data. In order to clean all the data, source and build data produced during the ISO build, the simplest way is to delete the docker volume after having exited the toolchain container. 


To exit the toolchain container use the command from within the command line of the running container:

     exit

Run the following command to remove the data.

    docker volume rm build_data

## License
APACHE 2.0