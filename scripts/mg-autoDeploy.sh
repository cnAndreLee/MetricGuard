#!/bin/bash
set -e

yes_flag=0

while getopts ":y" opt; do
    case $opt in
        y)
            yes_flag=1
            ;;
    esac
done

if [[ "$yes_flag" -eq 0 ]];then
    echo "Not allowed run without -y"
    exit 1
fi

cd $(dirname $0)

update=0
if [[ -e ~/shell/mg-release/mg_client.rc ]];then
    update=1
    bash ~/shell/mg-release/stop_client.sh
    bash ~/shell/mg-release/stop_server.sh
    sleep 2
    cp ~/shell/mg-release/mg_client.rc /tmp/mg_client.rc.bak
    sed -i 's/127.0.0.1:2999/localhost/g' /tmp//mg_client.rc.bak
    rm -rf ~/shell/mg-release
fi

if [[ -e ~/shell/client/mg_client.rc ]];then
    update=1
    bash ~/shell/mg-release/client/stop_client.sh
    bash ~/shell/mg-release/server/stop_server.sh
    sleep 2
    cp ~/shell/mg-release/client/mg_client.rc /tmp/mg_client.rc.bak
    rm -rf ~/shell/mg-release
fi

mkdir -p ~/shell/mg-release

cp -r ../mg-release/* ~/shell/mg-release/
cd ~/shell/mg-release

if [[ "$update" -eq 1 ]];then
    mv /tmp/mg_client.rc.bak ./private/mg_client.rc
    bash ~/shell/mg-release/server/start_server.sh
    bash ~/shell/mg-release/client/start_client.sh
    echo "已完成更新"
else
    cd private
    mv mg_client.rc.example mg_client.rc
    echo "已完成部署,请编辑配置文件后再进行启动"
fi

