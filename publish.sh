#!/bin/bash

#启动命令及参数
OutPath="./"
Target="teach"

if [ "$1" = "" ]; then
    echo -e "\033[0;31m 未输入环境参数 \033[0m  \033[0;34m {test|online} \033[0m"
    exit 1
fi

if [ "$Target" = "" ]; then
    echo -e "\033[0;31m 未填写编译目标 \033[0m"
    exit 1
fi

echo -e "target:$Target\033[0;33m build... \033[0m"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $Target

if [ x"$?" != x"0" ]; then
    echo -e "\033[0;31m 编译失败 \033[0m"
    exit 1
fi

echo -e "build $OutPath$Target done,\033[0;33m ready upload... \033[0m"

function do_test_env()
{
	pscp -P 22 -pw yx,0615 $OutPath$Target shane@49.233.3.177:/home/shane/server/$Target

	if [ x"$?" == x"0" ]; then
		echo -e "publish $Target to environment[\033[0;36m test \033[0m]\033[0;32m finish \033[0m."
	else
		echo -e "publish $Target to environment[\033[0;36m test \033[0m]\033[0;31m failure \033[0m."
    fi
}

function do_online_env()
{
	scp $OutPath$Target shane@49.233.3.177:/home/shane/server/$Target

	if [ x"$?" == x"0" ]; then
		echo -e "publish $Target to environment[\033[0;34m online \033[0m]\033[0;32m finish \033[0m."
	else
		echo -e "publish $Target to environment[\033[0;34m online \033[0m]\033[0;31m failure \033[0m."
    fi
}

case $1 in
    online)
    do_online_env;;
    test)
    do_test_env;;
    *)

esac
