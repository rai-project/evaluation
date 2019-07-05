#!/usr/bin/env bash

DATABASE_ADDRESS=$1
MODELNAME=ResNet_v1_50
BATCHSIZE=2
OUTPUTFOLDER=output

go build main.go

# run the analysis
# ./main latency --database_address=$DATABASE_ADDRESS --model_name=$MODELNAME --batch_size=$BATCHSIZE --format=csv --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/latency.csv"

# ./main cuda_kernel --database_address=$DATABASE_ADDRESS --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/cuda_kernel.csv"

./main layer --database_address=$DATABASE_ADDRESS --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_layer --format=csv --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer.csv"

# ./main layer --database_address=$DATABASE_ADDRESS --model_name=$MODELNAME --batch_size=$BATCHSIZE --bar_plot --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/barplot.html"

# ./main layer --database_address=$DATABASE_ADDRESS --model_name=$MODELNAME --batch_size=$BATCHSIZE --box_plot --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/boxplot.html"
