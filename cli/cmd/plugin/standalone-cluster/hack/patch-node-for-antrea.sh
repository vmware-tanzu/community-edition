#!/bin/bash

peerIdx=$(docker exec $1 ip link | grep eth0 | awk -F[@:] '{ print $3 }' | cut -c 3-)
peerName=$(docker run --net=host antrea/ethtool:latest ip link | grep ^"$peerIdx": | awk -F[:@] '{ print $2 }' | cut -c 2-)
docker run --net=host --privileged antrea/ethtool:latest ethtool -K "$peerName" tx off
docker exec "$1" sysctl -w net.ipv4.conf.all.route_localnet=1