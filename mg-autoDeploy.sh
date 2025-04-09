#!/bin/bash

mkdir -p ~/shell 

set -e
mv mg-release/ ~/shell/
cd ~/shell/mg-release
mv mg_client.rc.example mg_client.rc
vim mg_client.rc

