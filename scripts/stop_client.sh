#!/bin/bash

cd `dirname $0`
app="mg_client.sh"

count=$(ps -ef | grep $app | grep -v grep | wc -l)

if [[ $count = '0' ]];then
    echo "APP already stopped"
elif [[ $count = '1' ]];then
    ps -ef | grep $app | grep -v grep | awk '{print $2}' | xargs kill 
    echo "killed"
else 
    echo "Similar apps exist, Please handle "
fi

