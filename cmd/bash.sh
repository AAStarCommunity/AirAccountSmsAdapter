#!/bin/bash

exec="smsadapter"

echo "Killing and Starting"

pkill $exec
nohup ../$exec &

# 设置Git仓库URL和本地目录
repo_url="https://github.com/AAStarCommunity/AirAccountSmsAdapter.git"

# 设置主分支名称
branch="main"

while true; do
    git fetch origin "$branch"
    
    if [ "$(git rev-parse HEAD)" != "$(git rev-parse origin/$branch)" ]; then

        echo "Founding update, waiting to upgrade..."
        git pull origin "$branch"

        
        echo "Stoping..."
        pkill $exec
        
        echo "Updating"
        CGO_ENABLED=1 GOARCH=arm go build -o $exec main.go

        nohup ../$exec &
        
        echo "Updated"
    fi
    
    sleep 5
done
