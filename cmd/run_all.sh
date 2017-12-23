#!/bin/sh

go build main.go
./main latency --model_name all --overwrite -f json
./main layers --model_name all --overwrite -f json
./main layer_tree --model_name all --overwrite -f json
./main cuda_launch --model_name all --overwrite -f json
./main eventflow --model_name all --overwrite -f json

