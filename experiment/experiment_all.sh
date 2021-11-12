#!/bin/bash
experiments=("./experiment_cache_size.sh" "./experiment_associativity.sh" "./experiment_block_size.sh")
benchmarks=("blackscholes" "bodytrack" "fluidanimate")
protocols=("MESI" "MESIF" "Dragon")

for experiment in "${experiments[@]}"
do
    for protocol in "${protocols[@]}"
    do
        for benchmark in "${benchmarks[@]}"
        do
            $experiment $protocol $benchmark
        done
    done
done