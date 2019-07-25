#!/usr/bin/env bash

DATABASE_ADDRESS=$1
BATCHSIZE=$2
MODELNAME=MLPerf_ResNet50_v1.5
OUTPUTFOLDER=output
DATABASE_NAME=test

if [ ! -d $OUTPUTFOLDER ]; then
  mkdir $OUTPUTFOLDER
fi

if [ -f main ]; then
  rm main
fi

go build main.go


# ./main model info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/model"

# ./main model  info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/model.html"

echo "Start to run layer analysis"

./main layer info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer"

./main layer duration --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_duration_bar.html"

./main layer duration --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --box_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_duration_box.html"

./main layer memory --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_memory.html"

./main layer occurrence --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --pie_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_occurrence.html"

./main layer aggre_duration --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --pie_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/layer_aggre_duration.html"

echo "Start to run gpu analysis"

./main gpu_kernel info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel"

./main gpu_kernel name_aggre info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_name_aggre"

./main gpu_kernel model_aggre info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_model_aggre"

./main gpu_kernel layer_aggre info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv,table --output="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_aggre"

./main gpu_kernel layer_gpu_cpu info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_gpu_cpu.html"

./main gpu_kernel layer_flops info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_flops.html"

./main gpu_kernel layer_dram_read info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_dram_read.html"

./main gpu_kernel layer_dram_write info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --bar_plot --plot_path="$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/gpu_kernel_layer_dram_write.html"
