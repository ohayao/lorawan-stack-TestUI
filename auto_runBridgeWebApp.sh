#!/bin/bash

procnum=`ps -ef|grep "lorawan_test"|grep -v grep|wc -l`

if [ $procnum -eq 0 ]; then
    cd /root/project/lorawanTest
    rm -rf lorawan_test
    go build -o lorawan_test main.go lora-ttn.go sse.go
    nohup /root/project/lorawanTest/lorawan_test >> /root/project/lorawanTest/logs 2>&1 &
fi
