#!/bin/bash
for ((i=0; i<1000; i++))
do
	k=$(head -c 25 /dev/urandom | base64)
	v=$(head -c 40 /dev/urandom | base64)
	echo "${k},${v}"
done

