#!/bin/bash

CUR_DIR=$(cd "$(dirname "$0")";pwd)
GPU_MON="${CUR_DIR}/gpu-mon"
CFG_FILE="${CUR_DIR}/cfg.json"

if [ ! -x "$GPU_MON" ];then
    echo $GPU_MON is not exist
    exit 2
fi
if [ ! -f "$CFG_FILE" ]; then
    echo $CFG_FILE is not exist
    exit 3
fi
$GPU_MON -c $CFG_FILE -o