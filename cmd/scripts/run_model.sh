#!/bin/sh

cd ..
go build main.go
# ./main latency --model_name ResNet50 --model_version=1.0 --database_address=52.91.209.88 --database_name resnet50_1_0
./main latency --model_name ShuffleNet_Caffe2 --model_version=1.0 --database_address=192.17.100.10 --database_name shufflenet_model_trace
# ./main all --model_name SphereFace --model_version=1.0 --database_address=52.91.209.88 --database_name sphereface_v1_0_docker --hostname=abduld-nuc
# ./main latency --database_address=minsky1-1.csl.illinois.edu --database_name resnet_50_v1_0_model_trace --model_name ResNet50 --model_version=1.0
# ./main latency --append --model_name all --database_name resnet_50_1_0 &
# ./main layers	--append --model_name all -f json --database_name carml_full_trace &
#./main cuda_launch	--apend --model_name all -f json --database_name carml_full_trace &
#./main eventflow	--append --model_name all -f json --database_name carml_full_trace &
wait

