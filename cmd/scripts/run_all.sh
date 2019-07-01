#!/bin/sh

cd ..
go build main.go
./main latency --append --model_name all -f json --database_name carml_model_trace &
./main layers	--append --model_name all -f json --database_name carml_full_trace &
#./main cuda_launch	--apend --model_name all -f json --database_name carml_full_trace &
#./main eventflow	--append --model_name all -f json --database_name carml_full_trace &
wait

