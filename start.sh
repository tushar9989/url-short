#!/bin/bash
# Demo script to start read, write and auth servers

function changePortInConfig {
    if [ ! -f $2 ]; then
        echo "\"$2\" not found!"
        return
    fi
    sed -i "s|\(\"Port\": \)[0-9]*|\1$1|" $2
}

function changeServerIdInConfig {
    if [ ! -f $2 ]; then
        echo "\"$2\" not found!"
        return
    fi
    sed -i "s|\(\"ServerID\": \"\)[0-9]*|\1$1|" $2
}

cd auth-server
changePortInConfig 8001 "config.json"
go run main.go &
sleep 1

cd ../read-server/cmd/server
for (( i=8002; i<=8004; i++ ))
do
    changePortInConfig $i "../../config.json"
    go run main.go &
    sleep 1
done

cd ../../../write-server/cmd/server
for (( i=8005; i<=8006; i++ ))
do
    changePortInConfig $i "../../config.json"
    changeServerIdInConfig $((i-8005)) "../../config.json"
    go run main.go &
    sleep 1
done

wait