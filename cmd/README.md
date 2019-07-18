# Evaluation Commands

```go build main.go```

Run `./main -h` for help.

To get help on a subcommand, accuracy for example, run `go run main.go accuracy -h`.


## Model

* Model information across different batch sizes

   ```./main model info --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --format=csv```

## Layer

* Layer information

  ```./main layer info --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --format=csv```

* Layer duration

  ```./main layer duration --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --bar_plot```

* Layer duration variance

  ```./main layer duration --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --box_plot```

* Layer memory

  ```./main layer memory --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --bar_plot```

* Layer occurrence

  ```./main layer occurrence --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --pie_plot```

* Layer aggregated duration based on operator type

  ```./main layer aggre_duration --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --pie_plot```

* Layer theoretical flops calculation using the layer operator type and shape

  TODO

## GPU

* GPU kernel information

  ```./main gpu_kernel info --database_name=$DATABASE_NAME --database_address=$DATABASE_ADDRESS --model_name=$MODEL_NAME --output=$OUTPUTFILE --batch_size=$BATCH_SIZE --format=csv```

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

* Layer GPU vs CPU time

  Use the information from  ```gpu_kernel layer_aggre```

* GPU kernel roofline analysis

  Use the information from  ```gpu_kernel info```

* Layer roofline analysis

  Use the information from  ```gpu_kernel layer_aggre```

* Model roofline analysis

  Use the information from  ```gpu_kernel model_aggre```
