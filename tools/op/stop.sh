#!/bin/sh

echo "Stopping service..."

cd ../../
process=`pwd | xargs -i basename {}`


PIDS=`ps -ef | grep -E "./${process}"  | grep -v grep | awk '{print $2}'`

for MAIN_PID in $PIDS 
do
        kill -9 $MAIN_PID
done

PIDS=`ps -ef | grep -E "./${process}"  | grep -v grep | awk '{print $2}'`
for MAIN_PID in $PIDS 
do
        echo "Can't stop."
done

if [ $? != "0" ]
then
echo "Can't stop."
exit 1
else
echo "Service stoped."
exit 0
fi
