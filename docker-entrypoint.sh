#!/bin/sh

echo "*** docker-entorypoint.sh Start... ***"
cname=`cat ./cname`

rm ./$cname
go build -o $cname
exec ./$cname
