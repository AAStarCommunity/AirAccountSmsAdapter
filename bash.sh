#!/bin/bash

# 设置Git仓库URL和本地目录
repo_url="https://github.com/AAStarCommunity/AirAccountSmsAdapter.git"
local_dir="/src"

# 设置主分支名称
branch="main"
exec="smsadapter"

# 进入本地目录
cd "$local_dir"

while true; do
    git fetch origin "$branch"
    
    if [ "$(git rev-parse HEAD)" != "$(git rev-parse origin/$branch)" ]; then

        pkill $exec

        echo "Updating and restarting the program..."

        git pull origin "$branch"

        CGO_ENABLED=1 GOARCH=arm go build -o $exec main.go

        ./$exec
    fi
    
    sleep 60
done
