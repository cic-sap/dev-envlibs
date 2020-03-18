#!/usr/bin/env bash
set -ex


A=$(curl http://envlibs)
B="hello envlibs!"

echo "A:$A"
echo "B:$B"

#compare strings
if [ "$A" == "$B" ];then
exit 0
fi

exit 1
