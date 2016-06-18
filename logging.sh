#!/bin/bash

url=$1
bucket=${url#h*[e][r][/]}
bucket=${bucket%/*}
echo $bucket
filename=${bucket##**/}
datetime=$(date +"%m_%d_%Y-%H_%M_%S")
filename=$filename-$datetime
#echo $filename
mkdir $filename
cd $filename
gslink=gs://$bucket/*
echo $gslink
gsutil -m cp -r $gslink .
cd ..
go run rdjunit.go $filename > $filename/output.txt
