# Evaluation Commands

## Availabel Commands and Flags

Run

```
go run main.go -h
```

for help.

```
Get evaluation information from MLModelScope

Usage:
evaluation [command]

Available Commands:
accuracy Get accuracy summary from MLModelScope
all Get all evaluation information from MLModelScope
cuda_launch Get evaluation kernel launch information from MLModelScope
duration Get evaluation duration summary from MLModelScope
eventflow Get evaluation trace in event_flow format from MLModelScope
help Help about any command
latency Get evaluation latency or throughput information from MLModelScope
layer_tree Get evaluation layer tree information from MLModelScope
layers Get evaluation layer information from MLModelScope

Flags:
--append append the output
--arch string architecture of machine to filter
--batch_size int the batch size to filter
--database_address string address of the database
--database_name string name of the database to query
-f, --format string print format to use (default "table")
--framework_name string frameworkName
--framework_version string frameworkVersion
-h, --help help for evaluation
--hostname string hostname of machine to filter
--limit int limit the evaluations (default -1)
--model_name string modelName (default "BVLC-AlexNet")
--model_version string modelVersion (default "1.0")
--no_header show header labels for output
-o, --output string output file name
--overwrite if the file or directory exists, then they get deleted

Use "evaluation [command] --help" for more information about a command.

```

To get help on a subcommand, accuracy for example, run

```
go run main.go accuracy -h
```

## accuracy

Get accuracy summary from MLModelScope

## all

Get all evaluation information from MLModelScope
##cuda_launch
Get evaluation kernel launch information from MLModelScope

## duration

Get evaluation duration summary from MLModelScope

## eventflow

Get evaluation trace in event_flow format from MLModelScope

## latency

Get evaluation latency or throughput information from MLModelScope

## layer_tree

Get evaluation layer tree information from MLModelScope

## layers

Get evaluation layer information from MLModelScope
