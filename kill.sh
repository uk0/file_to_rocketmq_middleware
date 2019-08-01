#!/bin/sh
BASE=$(cd `dirname $0`; pwd)

cd ${BASE} && kill -9 `ps -ef | grep $(cat command.pid) | awk 'NR>1{print$2}'`