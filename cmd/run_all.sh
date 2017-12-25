#!/bin/sh

#go build main.go
./main latency	--append --model_name all -f json --database_name carml_step_trace &
./main layers	--append --model_name all -f json --database_name carml_full_trace &
./main layer_tree	--append --model_name all -f json --database_name carml_full_trace &
#./main cuda_launch	--apend --model_name all -f json --database_name carml_full_trace &
#./main eventflow	--append --model_name all -f json --database_name carml_full_trace &
wait

