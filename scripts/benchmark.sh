#!/bin/bash
sudo sysctl vm.drop_caches=3 && python main.py --path $1 && sudo sysctl vm.drop_caches=3 && go run main.go -path $1