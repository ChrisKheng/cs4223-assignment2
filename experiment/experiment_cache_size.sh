#!/bin/bash
# Run this script in the experiment directory
protocol=$1
benchmark_name=$2

if [[ $# -ne 2 ]]
then
    echo "Usage: experiment_cache_size.sh <protocol> <benchmark_name>"
    echo "protocol: MESI or Dragon"
    echo "benchmark name: blackscholes, bodytrack, or fluidanimate"
    exit 1
fi

output_dir="../experiment/${protocol}/cache_size/${benchmark_name}"
benchmark_dir="../benchmarks/${benchmark_name}_four/${benchmark_name}"

mkdir -p $output_dir

cd ../coherence
cache_size=64
for i in {0..10}
do
    echo "Running iter ${i}, cache size: ${cache_size}"
    ./coherence $protocol $benchmark_dir $cache_size 2 32 1> ${output_dir}/${cache_size}.out 2> ${output_dir}/${cache_size}.err
    (( cache_size *= 2 ))
done
