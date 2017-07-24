#!/bin/sh
cd $(dirname $0)
LD_PRELOAD=$PWD/libxd.so ./main

