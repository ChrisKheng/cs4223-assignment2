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

output_dir="../experiment/${protocol}/associativity/${benchmark_name}"
benchmark_dir="../benchmarks/${benchmark_name}_four/${benchmark_name}"

mkdir -p $output_dir

cd ../coherence
associativity=1
for i in {0..7}
do
    echo "Running iter ${i}, associativity: ${associativity}"
    ./coherence $protocol $benchmark_dir 4096 $associativity 32 1> ${output_dir}/${associativity}.out 2> ${output_dir}/${associativity}.err
    (( associativity *= 2 ))
done
