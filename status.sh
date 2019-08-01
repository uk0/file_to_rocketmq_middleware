#!/bin/sh
BASE=$(cd `dirname $0`; pwd)

cd ${BASE} && tail -f running.log  | grep SendMessageResultQueueOffset
