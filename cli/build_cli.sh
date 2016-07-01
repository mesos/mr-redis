#!/bin/bash


echo "BUILDING WINDOWS CLI"
echo "***************************************************"
GOOS=windows godep go build -v -o mrr_win_amd64 .
echo "BUILDING MAC CLI"
echo "***************************************************"
GOOS=darwin godep go build -v -o mrr_darwin_amd64 .
echo "BUILDING LINUX CLI"
echo "***************************************************"
GOOS=linux godep go build -v -o mrr_linux_amd64 .

