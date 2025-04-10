#!/bin/bash

cd $(dirname $0)

app="MetricGuard"

rm -rf ./mg-release
rm -f mg-release.tar.gz
mkdir -p mg-release/server/logs
mkdir -p mg-release/client/logs
mkdir -p mg-release/private

cp -r ./client/* ./mg-release/client/
cp ./client/mg_client.rc.example ./mg-release/private/
cp ./private/mg_server.rc ./mg-release/private/
cp ./scripts/start_server.sh ./mg-release/server/
cp ./scripts/stop_server.sh ./mg-release/server/
cp ./scripts/start_client.sh ./mg-release/client/
cp ./scripts/stop_client.sh ./mg-release/client/

# cp ./scripts/mg-autoDeploy.sh ./mg-release/

cd server
ARCH="arm64"
OS="linux"
CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -o ../mg-release/server/$app"_"$ARCH

ARCH="amd64"
OS="linux"
CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -o ../mg-release/server/$app"_"$ARCH

cd ../
tar -czpf mg-release.tar.gz mg-release scripts/mg-autoDeploy.sh

