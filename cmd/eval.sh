#!/usr/bin/env bash

DATABASE_ADDRESS=$1
BATCHSIZE=$2
MODELNAME=MLPerf_ResNet50_v1.5
OUTPUTFOLDER=output
DATABASE_NAME=carml

if [ ! -d $OUTPUTFOLDER ]; then
  mkdir $OUTPUTFOLDER
fi

if [ -f main ]; then
  rm main
fi

go build main.go

./main model info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --sort_output --format=csv,table --plot_all --output="$OUTPUTFOLDER/$MODELNAME/model_info"

echo "Start to run layer analysis"

# ./main layer info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --format=csv,table --plot_all --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_info"

# ./main layer aggre_info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --format=csv,table --plot_all --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_aggre_info"

# echo "Start to run gpu analysis"

# ./main gpu_kernel info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_info"

# ./main gpu_kernel name_aggre info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_name_aggre"

# ./main gpu_kernel model_aggre info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_model_aggre"

# ./main gpu_kernel layer_aggre info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_aggre"

# ./main gpu_kernel layer_gpu_cpu --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_gpu_cpu.html"

# ./main gpu_kernel layer_flops --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_flops.html"

# ./main gpu_kernel layer_dram_read --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_dram_read.html"

# ./main gpu_kernel layer_dram_write --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_dram_write.html"
