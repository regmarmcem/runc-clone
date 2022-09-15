#!/bin/bash
PATH=$PATH:$(pwd)
go build -o runc-clone && sudo ./runc-clone run --mount ./bundle/rootfs --uid 0 --debug --command "ls"
