#!/bin/sh
set -ex

docker run -d \
    --restart always \
    --name mgmt \
    --net host \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /boot:/boot_host \
    -v /config:/config_host \
    -v /var/config:/varconfig_host \
    --cap-add NET_ADMIN \
    fouwels.azurecr.io/fouwels/mgmt:dev
