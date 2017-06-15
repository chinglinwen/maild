#!/bin/sh

rm *.gz
rm maild maild.exe
GOOS=linux GOARCH=386 go build
GOARCH=386 go build
tar -czf maild_win386.tar.gz maild.exe
tar -czf maild_linux386.tar.gz maild
echo done