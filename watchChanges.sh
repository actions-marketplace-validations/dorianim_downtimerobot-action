#!/bin/bash


sync() {
    go run . frontend
    #rsync -avzp --chown 0:0 ./usr root@$VM_IP:/
    #rsync -avzp --chown 0:0 ./etc root@$VM_IP:/
}

sync

inotifywait -r -m -e close_write --format '%w%f' ./internal | while read MODFILE
do
    echo need to rsync ${MODFILE%/*} ...
    sync
done