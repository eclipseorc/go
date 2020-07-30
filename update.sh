#!/bin/bash

#应用拷贝源
SrcName=../teach
#需要启动的应用名
AppName=teach_8007
#启动命令及参数
StartComment="./$AppName"

if [ "$1" = "" ];
then
    echo -e "\033[0;31m 未输入动作名 \033[0m  \033[0;34m {start|stop|restart|status} \033[0m"
    exit 1
fi

if [ "$AppName" = "" ];
then
    echo -e "\033[0;31m 未填写程序名 \033[0m"
    exit 1
fi

function do_start()
{
	PID=`pidof $AppName` 
    # PID=`ps -ef |grep $AppName|grep -v grep|awk '{print $2}'`
	
	if [ x"$PID" != x"" ]; then
	    echo -e "$AppName is \033[0;33m running...\033[0m"
	else
		if [ "$SrcName" != "$AppName" ]; then
			echo -e "\033[0;33m 更新程序... \033[0m"
			cp $SrcName $AppName
		fi

		# nohup $StartComment >/dev/null 2>&1 &
		nohup $StartComment &

		sleep 2
		PID=`pidof $AppName` 
		if [ x"$PID" != x"" ]; then
			echo -e "Start $AppName \033[0;32m success\033[0m."
		else
			echo -e "Start $AppName \033[0;31m failure, please check the logs\033[0m."
		fi
	fi
}

function do_stop()
{
    echo -e "$AppName \033[0;36m stopping...\033[0m"
	
	PID=""
	query(){
		PID=`pidof $AppName` 
		# PID=`ps -ef |grep $AppName|grep -v grep|awk '{print $2}'`
	}
 
	query
	if [ x"$PID" != x"" ]; then

		# kill -TERM $PID
		kill -3 $PID && sleep 2 && kill -TERM $PID

		echo -e "$AppName (pid:$PID) \033[0;36m exiting...\033[0m"
		while [ x"$PID" != x"" ]
		do
			sleep 1
			query
		done
		echo -e "$AppName \033[0;32m exited\033[0m."
	else
		echo -e "$AppName \033[0;33m already stopped\033[0m."
	fi
}

function do_restart()
{
    do_stop
    sleep 1
	backfile
    do_start
}

function do_status()
{
	PID=`pidof $AppName` 
    # PID=`ps -ef |grep $AppName|grep -v grep|awk '{print $2}'`
	
	if [ x"$PID" != x"" ]; then
        echo -e "$AppName is \033[0;36m running...\033[0m"
    else
        echo -e "$AppName is \033[0;33m not running...\033[0m"
    fi
}

function backfile()
{
	BackSuffix='.bak'
	Modify=`stat -c %Y $AppName | date '+_%Y-%m-%d_%H:%M:%S'`
	cp $AppName ${SrcName##*/}${Modify}${BackSuffix}
	echo -e "Back $AppName to ${SrcName##*/}${Modify}${BackSuffix} \033[0;32m finish\033[0m."
}
case $1 in
    start)
    do_start;;
    stop)
    do_stop;;
    restart)
    do_restart;;
    status)
    do_status;;
    *)
 
esac