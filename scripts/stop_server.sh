#!/bin/bash

cd `dirname $0`
app="MetricGuard"

count=$(ps -ef | grep $app | grep -v grep | wc -l)

if [[ $count = '0' ]];then
    echo "APP already stopped"
elif [[ $count = '1' ]];then
    ps -ef | grep $app | grep -v grep | awk '{print $2}' | xargs kill 
    echo "Killed"
else 
    echo "Similar apps exist, Please handle "
fi

