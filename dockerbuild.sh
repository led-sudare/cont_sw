#!/bin/sh
cname=`cat ./cname`
docker build ./ -t $cname
docker container stop $cname
docker container rm $cname
docker run -itd --init --name $cname -v `pwd`:/work/ -p 8002:8002 $cname