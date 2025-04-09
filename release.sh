#!/bin/bash

cd $(dirname $0)

app="MetricGuard"

rm -rf ./mg-release
rm -f mg-release.tar.gz
mkdir -p mg-release/logs

cp ./client/mg_client.sh ./mg-release/
cp ./private/mg_server.rc ./mg-release/
cp ./client/mg_client.rc.example ./mg-release/

cp ./scripts/start_server.sh ./mg-release/
cp ./scripts/stop_server.sh ./mg-release/
cp ./scripts/start_client.sh ./mg-release/
cp ./scripts/stop_client.sh ./mg-release/

cd server
ARCH="arm64"
OS="linux"
CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -o ../mg-release/$app"_"$ARCH

ARCH="amd64"
OS="linux"
CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -o ../mg-release/$app"_"$ARCH

cd ../
tar -czpf mg-release.tar.gz mg-release mg-autoDeploy.sh mg-autoUpdate.sh
