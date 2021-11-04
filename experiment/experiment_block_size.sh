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

output_dir="../experiment/block_size/${benchmark_name}"
benchmark_dir="../benchmarks/${benchmark_name}_four/${benchmark_name}"

mkdir -p $output_dir

cd ../coherence
block_size=4
for i in {0..9}
do
    echo "Running iter ${i}, block size: ${block_size}"
    ./coherence $protocol $benchmark_dir 4096 2 ${block_size} 1> ${output_dir}/${block_size}.out 2> ${output_dir}/${block_size}.err
    (( block_size *= 2 ))
done
