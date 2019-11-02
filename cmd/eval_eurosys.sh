#!/usr/bin/env bash

DATABASE_ADDRESS=$1
BATCHSIZE=$2
MODELNAME=$3
OUTPUTFOLDER=$4
DATABASE_NAME=$5


if [ ! -d $OUTPUTFOLDER ]; then
  mkdir -p $OUTPUTFOLDER
fi


echo Start to run layer analysis

./main model info --database_address=$DATABASE_ADDRESS --database_name=$DATABASE_NAME --model_name=$MODELNAME --batch_size=$BATCHSIZE --sort_output --format=csv --plot_all --output=$OUTPUTFOLDER/$MODELNAME/$BATCHSIZE/model_info

