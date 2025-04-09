#!/bin/bash

cd `dirname $0`

nohup bash -c '
while true; do
    nohup bash mg_client.sh -y >> ./logs/mg_client.log 2>&1 &
    sleep 300
done
' > /dev/null 2>&1 &
