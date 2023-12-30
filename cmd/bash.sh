#!/bin/bash

local_dir=".."
exec="smsadapter"

echo "Starting"
$local_dir/$exec

# 设置Git仓库URL和本地目录
repo_url="https://github.com/AAStarCommunity/AirAccountSmsAdapter.git"

# 设置主分支名称
branch="main"

while true; do
    git fetch origin "$branch"
    
    if [ "$(git rev-parse HEAD)" != "$(git rev-parse origin/$branch)" ]; then

        pkill $exec

        echo "Updating and restarting the program..."

        git pull origin "$branch"

        CGO_ENABLED=1 GOARCH=arm go build -o $local_dir/$exec $local_dir/main.go

        $local_dir/$exec
    fi
    
    sleep 60
done
