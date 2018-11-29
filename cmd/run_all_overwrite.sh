#!/bin/sh

go build main.go
./main latency --model_name all --overwrite -f json --database_name carml_model_trace &
./main layers --model_name all --overwrite -f json --database_name carml_full_trace &
./main layer_tree --model_name all --overwrite -f json --database_name carml_full_trace &
#./main cuda_launch --model_name all --overwrite -f json --database_name carml_full_trace &
#./main eventflow --model_name all --overwrite -f json --database_name carml_full_trace &
wait

