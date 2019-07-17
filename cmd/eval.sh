#!/usr/bin/env bash

DATABASE_ADDRESS=$1
MODELNAME=ResNet_v1_50
BATCHSIZE=2
NUMPREDS=3
DUPLICATE_INPUT=$(($NUMPREDS * $BATCHSIZE))
OUTPUTFOLDER=output
DATABASE_NAME=carml

if [ ! -d $OUTPUTFOLDER ]; then
  mkdir $OUTPUTFOLDER
fi

# run the analysis
# go run main.go model info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --format=csv --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/model.csv"
# go run main.go model info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/model.tbl"

# go run main.go layer info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --format=csv --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer.csv"
# go run main.go layer info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer.tbl"

# go run main.go layer latency --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_latency.csv"

# go run main.go layer duration --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --pie_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_duration.html"

# go run main.go layer occurrence --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --pie_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_occurrence.html"

# go run main.go layer memory --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_memory.html"

# go run main.go cuda_kernel info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/cuda_kernel.csv"
# go run main.go cuda_kernel info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/cuda_kernel.tbl"

go run main.go layer cuda_kernel --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_cuda_kernel.csv"
go run main.go layer cuda_kernel --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_cuda_kernel.tbl"
