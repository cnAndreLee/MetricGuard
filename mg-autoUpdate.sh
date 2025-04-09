#!/bin/bash

cd $(dirname $0)

bash ~/shell/mg-release/stop_server.sh
sleep 2

cp ./mg-release/MetricGuard_amd64 ~/shell/mg-release/
cp ./mg-release/MetricGuard_arm64 ~/shell/mg-release/
cp ./mg-release/mg_server.rc ~/shell/mg-release/

cp ./mg-release/mg_client.sh ~/shell/mg-release/
rm -f ~/shell/mg-release/logs/mg_client.log

bash ~/shell/mg-release/start_server.sh

