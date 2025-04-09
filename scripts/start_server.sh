#!/bin/bash

cd `dirname $0`

source ./mg_server.rc

# 获取并标准化架构信息
ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')

# 判断架构类型
if [[ $ARCH =~ x86_64|amd64|i[3456]86 ]]; then
    echo "系统架构: x86/x64 (检测到: $ARCH)"
    nohup ./MetricGuard_amd64 >> ./logs/MetricGuard.log 2>&1 &
elif [[ $ARCH =~ arm|aarch64 ]]; then
    echo "系统架构: ARM (检测到: $ARCH)"
    nohup ./MetricGuard_arm64 >> ./logs/MetricGuard.log 2>&1 &
    # ARM相关操作
elif [[ $ARCH =~ ppc|powerpc ]]; then
    echo "系统架构: PowerPC (检测到: $ARCH)"
    # PowerPC相关操作
elif [[ $ARCH =~ mips ]]; then
    echo "系统架构: MIPS (检测到: $ARCH)"
    # MIPS相关操作
else
    echo "系统架构: 其他/未知 (检测到: $ARCH)"
fi

