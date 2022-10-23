#!/bin/bash
sudo sysctl vm.drop_caches=3 && python main.py  && sudo sysctl vm.drop_caches=3 && go run main.go