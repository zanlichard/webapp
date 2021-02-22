#!/bin/sh

cd ../../
process=`pwd | xargs -i basename {}`
ps -ef 2>/dev/null | awk '{ if( FNR == 1 ) printf "%s\n", $0;}'
ps -ef 2>/dev/null | grep "./bin/${process}" | grep -v grep | awk '{printf "%s\n", $0;}'

exit 0
