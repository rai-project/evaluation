# Evaluation Commands

```go build main.go```

Run `./main -h` for help.

To get help on a subcommand, accuracy for example, run `go run main.go accuracy -h`.


## Model

* model information across different batch sizes

   ```./main model info --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --format=csv```

## Layer

* layer information

  ```./main layer info --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --format=csv```

* layer duration

  ```./main layer duration --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --bar_plot```

* layer duration variance

  ```./main layer duration --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --box_plot```

* layer memory

  ```./main layer memory --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --bar_plot```

* layer occurrence

  ```./main layer occurrence --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --pie_plot```

* layer aggregated duration based on operator type

  ```./main layer aggre_duration --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --pie_plot```

* layer theoretical flops calculation using the layer operator type and shape

  TODO

## GPU

* GPU kernel information

  ```./main gpu_kernel info --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --format=csv```

  add ```--to
* GPU kernel information aggregated within each layer

  ```./main gpu_kernel layer_aggre --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --format=csv```

* GPU kernel information aggregated within the model

  ```./main gpu_kernel model_aggre --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --format=csv```

* Total flops of GPU kernels per layer

  ```./main gpu_kernel layer_flops --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --bar_plot```

* Total dram read of GPU kernels per layer

  ```./main gpu_kernel layer_dram_read --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --bar_plot```

* Total dram write of GPU kernels per layer

  ```./main gpu_kernel layer_dram_write --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --bar_plot```

* layer GPU vs CPU time

  TODO

* GPU kernel roofline analysis

  Use the information from  ```gpu_kernel info```

* layer roofline analysis

  Use the information from  ```gpu_kernel layer_aggre```

* model roofline analysis

  Use the information from  ```gpu_kernel model_aggre```
