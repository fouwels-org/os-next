#!/bin/bash

for module in $(cat /proc/modules | cut -d ' ' -f1 ) 
do
  modprobe -D "$module" | cut -d ' ' -f2 | sed 's/.*kernel/\/kernel/' >> ./tmp.dat
done

sort -u ./tmp.dat