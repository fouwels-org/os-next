#!/bin/sh
set -ex

docker run -d \
    --restart always \
    --name management \
    --net host \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /boot:/host/boot \
    -v /config:/host/config \
    -v /var/config:/host/var/config \
    --cap-add NET_ADMIN \
    registry2.lagoni.co.uk/belcan-as/alpine-management:1.0.0
