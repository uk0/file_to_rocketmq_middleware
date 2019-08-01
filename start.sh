#!/bin/sh
ulimit -n 655350

BASE=$(cd `dirname $0`; pwd)

cd ${BASE} && nohup ./FileToMQSender_linux > /dev/null 2>&1 &echo $!> command.pid
