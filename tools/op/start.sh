#!/bin/sh
#start_proc.sh

ulimit -c unlimited
echo "core-%e-%p-%t" > /proc/sys/kernel/core_pattern


start_single()
{
        now_dir=`pwd`
        cd ../../
        module=$(pwd | xargs -i basename {})
        count=$(ps -eo cmd  | awk '{print $1 }' | grep "^./bin/${module}$"| grep -v grep| wc -l)
        if [ $count -lt 1 ]
        then
               echo "[`date +'%Y-%m-%d %T'`] process ${module} number:$count, fork it!"
               echo `pwd`
               nohup ./bin/${module}  > ./logs/nohup.err 2>&1 &

        fi
        cd "$now_dir" > /dev/null 2>&1

}

start_single 
./p.sh
