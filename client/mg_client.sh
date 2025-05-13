#!/bin/bash

# 监控指标收集及告警v1.0

# 设置英文环境，以免脚本执行出错
export LANG=en_US.UTF-8

cd `dirname $0`

# 判断有没有加-y参数,加了-y参数运行才会发送告警
yes_flag=0
debug_flag=0

while getopts "yd" opt; do
    case $opt in
        y)
            yes_flag=1
            ;;
        d)
            debug_flag=1
            ;;
    esac
done

echo "debug_flag=$debug_flag"
echo "yes_flag=$yes_flag"

# 报警日志记录，根据此文件查询历史报警，判断发送逻辑
ALERT_LOG="./logs/alerts.log"
#配置文件，配置本机名称、需要检测的端口及短信接口地址

set -e
source ./mg_client.rc
set +e

echo "本机名称:$CMI_HOST"
echo "服务端口:$CMI_PORTS"
echo "短信接口:$CMI_SERVER_ENDPOINT"

# 网卡名称（目前是自动找第一个物理网卡，也可以指定）
interface=$(ls /sys/class/net -l | grep pci | awk '{print $9}' | head -n 1)
localIP=$(ip a show dev $interface | grep 'inet ' | awk '{print $2}' | cut -d/ -f1 )

function cpu_usage {

    # 第一次读取
    read cpu user nice system idle iowait irq softirq steal guest guest_nice < <(grep '^cpu ' /proc/stat)
    total1=$((user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice))
    user1=$user
    nice1=$nice
    system1=$system
    idle1=$idle
    iowait1=$iowait
    irq1=$irq
    softirq1=$softirq

    sleep 1

    # 第二次读取
    read cpu user nice system idle iowait irq softirq steal guest guest_nice < <(grep '^cpu ' /proc/stat)
    total2=$((user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice))
    user2=$user
    nice2=$nice
    system2=$system
    idle2=$idle
    iowait2=$iowait
    irq2=$irq
    softirq2=$softirq

    # 计算差值
    total=$((total2 - total1))
    user=$((user2 - user1))
    nice=$((nice2 - nice1))
    system=$((system2 - system1))
    idle=$((idle2 - idle1))
    iowait=$((iowait2 - iowait1))
    irq=$((irq2 - irq1))
    softirq=$((softirq2 - softirq1))

    # 计算各项占比
    cpu_percent=$((100 * ( total - idle ) / total))
    idle_percent=$((100 * idle / total))
    user_percent=$((100 * user / total))
    system_percent=$((100 * system / total))
    nice_percent=$((100 * nice / total))
    iowait_percent=$((100 * iowait / total))
    irq_percent=$((100 * irq / total))
    softirq_percent=$((100 * softirq / total))
    other_percent=$((100 * ( total - ( user + nice + system + idle + iowait + irq + softirq)) / total))
    
    metric_deal "cpu_usage" "CPU使用率" "$cpu_percent" "%" "$1" "$2" "$3"

    # 输出结果
    # echo "CPU 使用率(%): $cpu_percent"
    # echo "CPU 空闲时间占比: $idle_percent%"
    # echo "用户空间 CPU 使用率: $user_percent%"
    # echo "内核空间 CPU 使用率: $system_percent%"
    # echo "Nice 进程 CPU 使用率: $nice_percent%"
    # echo "iowait 状态占比: $iowait_percent%"
    # echo "CPU 中断时间占比: $irq_percent%"
    # echo "CPU 软中断时间占比: $softirq_percent%"
    # echo "其他 CPU 使用率: $other_percent%"
}

function cpu_load_avg1 {

    read cpu_load < <(echo "scale=2;`awk '{print $1}' /proc/loadavg` / `nproc`" | bc | awk '{printf "%.2f\n", $0}' )

    metric_deal "cpu_load_avg1" "CPU1分钟核心平均负载" "$cpu_load" "1" "$1" "$2" "$3"
}

function mem_usage {

    local MemAvailable=`cat /proc/meminfo | grep MemAvailable | awk '{print $2}'`;
    local MemTotal=`cat /proc/meminfo | grep MemTotal | awk '{print $2}'`;
    local mem_u=$((100 * (MemTotal - MemAvailable) / MemTotal));

    metric_deal "mem_usage" "内存使用率" "$mem_u" "%" "$1" "$2" "$3"
}

function total_open_files {
    
    local open_files=`lsof -u $USER | wc -l`

    metric_deal "total_open_files" "打开的文件数" "$open_files" "个" "$1" "$2" "$3"
}
function descriptors_limit {

    local descriptors_limit=$(ulimit -n)

    metric_deal "descriptors_limit" "最大描述符数" "$descriptors_limit" "" "$1" "$2" "$3"
}


function disk_used {
 
    df | grep mapper | while read -r line; do
        local mountPoint=$(echo $line | awk '{print $6}')
        local value=$(echo $line | awk '{print $5}' | sed 's/%$//')
        metric_deal "mount_point_${mountPoint}_used" "'${mountPoint}'挂载点空间使用率" "$value" "%" "$1" "$2" "$3"
    done

}

function disk_fs_rwstat {

    local fs=`df /home | tail -n 1 | awk '{print $1}'`
    local stat=$(grep $fs /proc/mounts | awk '{print $4}' | cut -f 1 -d ,)

    metric_deal "disk_fs_rwstat" "/home读写状态" "$stat" "" "$1" "$2" "$3"
}

function disk_inode_usage {

    df -i | grep mapper | while read -r line; do
        local mountPoint=$(echo $line | awk '{print $6}')
        local value=$(echo $line | awk '{print $5}' | sed 's/%$//')
        metric_deal "mount_point_${mountPoint}_inode_used" "'${mountPoint}'挂载点inode使用率" "$value" "%" "$1" "$2" "$3"
    done
}

function disk_io {

    fs=`df /home | tail -n 1 | awk '{print $1}'`
    fs_real_withpath=$(readlink $fs)
    fs_real=$(basename $fs_real_withpath)
    
    #第一次读取,/proc/diskstats
    read major minor name rio rmerge rsect ruse wio wmerge wsect wuse running use aveq < <(grep $fs_real /proc/diskstats | head -n 1)
    prev_rio=$rio
    prev_ruse=$ruse
    prev_wio=$wio
    prev_wuse=$wuse

    # 等待一秒
    sleep 1

    # 第二次读取/proc/diskstats
    read major minor name rio rmerge rsect ruse wio wmerge wsect wuse running use aveq < <(grep $fs_real /proc/diskstats | head -n 1)
    rio_diff=$((rio - prev_rio))
    ruse_diff=$((ruse - prev_ruse))
    wio_diff=$((wio - prev_wio))
    wuse_diff=$((wuse - prev_wuse))

    # 计算平均耗时
    if [ $rio_diff -gt 0 ]; then
        avg_ruse=$((ruse_diff / rio_diff))
    else
        avg_ruse=0
    fi

    if [ $wio_diff -gt 0 ]; then
        avg_wuse=$((wuse_diff / wio_diff))
    else
        avg_wuse=0
    fi

    metric_deal "disk_io_r" "读操作平均耗时" "$avg_ruse" "ms" "$1" "$2" "$3"
    metric_deal "disk_io_w" "写操作平均耗时" "$avg_wuse" "ms" "$1" "$2" "$3"
}

function net_drop_rate {

    # 网卡名称（目前是自动找第一个物理网卡，也可以指定）
    # interface="enp0s3"
    # interface=$(ls /sys/class/net -l | grep pci | awk '{print $9}' | head -n 1)

    # 第一次读取 /proc/net/dev
    read prev_rx_packets prev_rx_drops prev_tx_packets prev_tx_drops < <(awk -v iface="$interface" '$0 ~ iface ":" {print $3, $5, $11, $13}' /proc/net/dev)

    # 等待一秒
    sleep 1

    # 第二次读取 /proc/net/dev
    read curr_rx_packets curr_rx_drops curr_tx_packets curr_tx_drops < <(awk -v iface="$interface" '$0 ~ iface ":" {print $3, $5, $11, $13}' /proc/net/dev)

    # 计算差值
    rx_packets_diff=$((curr_rx_packets - prev_rx_packets))
    rx_drops_diff=$((curr_rx_drops - prev_rx_drops))
    tx_packets_diff=$((curr_tx_packets - prev_tx_packets))
    tx_drops_diff=$((curr_tx_drops - prev_tx_drops))

    # 计算丢包率
    if [ $rx_packets_diff -gt 0 ]; then
        rx_drop_rate=$(echo "scale=4; $rx_drops_diff / $rx_packets_diff * 100" | bc)
    else
        rx_drop_rate=0
    fi
    if [ $tx_packets_diff -gt 0 ]; then
        tx_drop_rate=$(echo "scale=4; $tx_drops_diff / $tx_packets_diff * 100" | bc)
    else
        tx_drop_rate=0
    fi

    # 输出结果
    # echo "Interface: $interface"
    # echo "Received packets (last second): $rx_packets_diff"
    # echo "Transmited packets (last second): $tx_packets_diff"
    # echo "RX Dropped packets (last second): $rx_drops_diff"
    # echo "TX Dropped packets (last second): $tx_drops_diff"

    metric_deal "rx_drop_rate" "接收丢包率" "$rx_drop_rate" "%" "$1" "$2" "$3"
    metric_deal "tx_drop_rate" "发送丢包率" "$tx_drop_rate" "%" "$1" "$2" "$3"
}

function net_tcp_time_wait {

    local count=$(ss -s | grep -oP 'timewait \K\d+')

    metric_deal "net_tcp_time_wait" "time-wait状态连接数" "$count" "个" "$1" "$2" "$3"
}

function net_tcp_close_wait {

    local operator=$1
    local threshold=$2
    local count=$(ss -t state close-wait | grep -c 'CLOSE-WAIT')

    metric_deal "net_tcp_close_wait" "close-wait状态连接数" "$count" "个" "$1" "$2" "$3"
}

function net_tcp_retrans_rate {

    # 第一次采样
    local retrans1=$(awk '/Tcp:/ {print $13}' /proc/net/snmp | tail -n 1)
    local seg_out1=$(awk '/Tcp:/ {print $12}' /proc/net/snmp | tail -n 1)

    # 等待 3 秒
    sleep 3

    # 第二次采样
    local retrans2=$(awk '/Tcp:/ {print $13}' /proc/net/snmp | tail -n 1)
    local seg_out2=$(awk '/Tcp:/ {print $12}' /proc/net/snmp | tail -n 1)

    # 计算差值
    local retrans=$((retrans2 - retrans1))
    local seg_out=$((seg_out2 - seg_out1))

    # 计算重传率
    if [ "$seg_out" -gt 0 ]; then
        retrans_rate=$(echo "scale=4; $retrans / $seg_out * 100" | bc)
    else
        retrans_rate=0
    fi

    metric_deal "net_tcp_retrans_rate" "TCP重传率" "$retrans_rate" "%" "$1" "$2" "$3"
}

function time_offset {

    local offset=""
    if chronyc tracking > /dev/null 2>&1 ;then 
        offset=$(chronyc tracking | grep "System" | awk '{print $4}')
    elif ntpq -pn > /dev/null 2>&1 ;then
        # offset=$(ntpq -pn | awk '/^\*/ {print $9}')
        offset=$(ntpq -pn | awk '/^\*/ {offset=$9; print (offset < 0) ? -offset : offset}')
    else
        echo "ERROR: 时间偏移量time_offset获取失败"
        return 1
    fi

    metric_deal "time_offset" "时间偏移量" "$offset" "seconds" "$1" "$2" "$3"
}

function check_port {

    local ports=($CMI_PORTS)
    local port_stat=""

    for port in "${ports[@]}"
    do 

        if nc -z $localIP $port ; then
            port_stat="UP"
        else 
            if nc -z 127.0.0.1 $port ; then
                port_stat="UP"
            else
                port_stat="DOWN"
            fi
        fi

        metric_deal "port_${port}_stat" "${port}端口状态" "$port_stat" "" "$1" "$2" "$3" 

    done
}

function EchoDebug {
    if [[ "$debug_flag" -eq 1 ]];then
        echo "$1"
    fi
}

# 收集并处理指标信息
function metric_deal {
    local metricName=$1
    local metricCnName=$2
    local metricValue=$3
    local metricUnit=$4
    local operator=$5
    local threshold5=$6
    local threshold4=$7

    local current_time=$(date +"%Y%m%d%H%M%S")  # 当前时间
    local current_unix_time=$(date -d "$(echo $current_time | sed -r 's/([0-9]{4})([0-9]{2})([0-9]{2})([0-9]{2})([0-9]{2})([0-9]{2})/\1-\2-\3 \4:\5:\6/')" "+%s")

    # 打印检测结果
    if [[ "$metricUnit" = "" ]];then
        echo "$current_time----${metricCnName}:${metricValue}"
    else
        echo "$current_time----${metricCnName}(${metricUnit}):${metricValue}"
    fi

    # 没有-y则退出
    if [[ $yes_flag -eq 0 ]]; then
        return 0
    fi

    # 报警类型 0:正常 4:四级 5:五级
    local AlertType
    case "$operator" in
        "<="|">=")
            if [[ $(echo "$metricValue $operator $threshold5" | bc) -eq 1 ]]; then
                EchoDebug "五级报警$metricName"
                AlertType=5
            elif [[ $(echo "$metricValue $operator $threshold4" | bc) -eq 1 ]]; then
                EchoDebug "四级报警$metricName"
                AlertType=4
            else
                EchoDebug "指标正常$metricName"
                AlertType=0
            fi
            ;;
        "!=")
            if [ "$metricValue" != "$threshold5" ]; then
                AlertType=5
            else
                AlertType=0
            fi
            ;;
        *)
            echo "不支持的operator"
            return 1
            ;;
    esac

    
    local last_alert_time

    local last_alert_type=0
    local needSend="N"

    # 读取上次指标报警或恢复的时间
    if [[ -f "$ALERT_LOG" ]]; then
        if [ $(grep $metricName $ALERT_LOG | tail -n 1 | wc -l) -eq "1" ];then 
            EchoDebug "存在指标记录$metricName"
            last_alert_time=$(grep $metricName $ALERT_LOG | tail -n 1 | cut -f1 -d#)
            last_alert_unix_time=$(date -d "$(echo $last_alert_time | sed -r 's/([0-9]{4})([0-9]{2})([0-9]{2})([0-9]{2})([0-9]{2})([0-9]{2})/\1-\2-\3 \4:\5:\6/')" "+%s")
            last_alert_type=$(grep $metricName $ALERT_LOG | tail -n 1 | cut -f3 -d#)
            local time_diff=$((current_unix_time - last_alert_unix_time))
            if ((time_diff < 86400)); then  # 24 小时
                EchoDebug "记录时间小于24小时$metricName"
                if [[ "$last_alert_type" = "$AlertType" ]];then
                    EchoDebug "alerttype相同,and in 24hours,needSend=N  $metricName"
                    needSend="N"
                elif [[ "$AlertType" = "5" ]];then
                    if [[ "$last_alert_type" = "5" ]];then
                        needSend="N"
                    else
                        needSend="Y"
                    fi
                elif [[ "$AlertType" = "4" ]];then
                    if [[ "$last_alert_type" = "0" ]];then
                        needSend="Y"
                    fi
                elif [[ "$AlertType" = "0" ]];then
                    if [[ "$last_alert_type" = "5" ]] || [[ "$last_alert_type" = "4" ]];then
                        needSend="Y"
                    fi
                fi 
            else 
                EchoDebug "记录时间大于24小时$metricName"
                if [[ "$AlertType" == "5" ]] || [[ "$AlertType" == "4" ]];then
                    EchoDebug "alerttype为告警$metricName"
                    needSend="Y"
                fi
                if [[ "$last_alert_type" == "5" ]] || [[ "$last_alert_type" == "4" ]];then
                    EchoDebug "超过24小时重复发送告警$metricName"
                    needSend="Y"
                fi
            fi
        else
            EchoDebug "不存在指标记录,4或5则发送 $metricName"
            # local noRecordFlag=true
            if [[ "$AlertType" == "5" ]] || [[ "$AlertType" == "4" ]];then
                needSend="Y"
            fi
        fi
    else
        EchoDebug "不存在指标记录且不存在日志文件,4或5则发送 $metricName"
        if [[ "$AlertType" == "5" ]] || [[ "$AlertType" == "4" ]];then
            needSend="Y"
        fi
    fi

    EchoDebug "needSend is $needSend , AlertType is $AlertType  $metricName"

    if [[ "$needSend" = "Y" ]];then

        local curlResult=0
        # 还需加入上一次报警类型
        set -x
        curl --unix-socket /tmp/metricguard.sock -X POST -H $'Content-Type: application/json' -d "{\"host\":\"$CMI_HOST\",\"ip\":\"$localIP\",\"type\":\"$AlertType\",\"last_type\":\"$last_alert_type\",\"metric_cn_name\": \"$metricCnName\",\"metric_value\": \"$metricValue\",\"metric_unit\": \"$metricUnit\",\"threshold5\": \"$threshold5\",\"threshold4\": \"$threshold4\"}" $CMI_SERVER_ENDPOINT
        local curlResult=$?
        set +x

        echo "$current_time#$metricName#$AlertType#$metricValue#$curlResult" >> "$ALERT_LOG"
    fi
}

cpu_usage '>=' 90 80 &
cpu_load_avg1 '>=' 0.9 0.7 &
mem_usage '>=' 90 80 &
total_open_files '>=' 60000 30000 &
descriptors_limit '<=' 30000 60000 &
disk_used '>=' 90 80 &
disk_fs_rwstat '!=' "rw" "rw" &
disk_inode_usage '>=' 90 80 &
disk_io '>=' 16 8 &
net_drop_rate '>=' 20 10 &
net_tcp_time_wait '>=' 1000 500 &
net_tcp_close_wait '>=' 30 20 &
net_tcp_retrans_rate '>=' 20 10 &
time_offset '>=' 2 1 &
check_port '!=' "UP" "UP" &

wait
