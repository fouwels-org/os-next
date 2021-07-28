#!/bin/bash

echo "" > modules.txt
for module in $(cat /proc/modules | cut -d ' ' -f1 ) 
do
  echo "$module" >> modules.txt
  modprobe -D "$module" | cut -d ' ' -f2 | sed 's/.*kernel/\/kernel/' >> modules.txt
done