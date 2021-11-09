#!/bin/bash
protocols=(MESI MESIF)
dimensions=(associativity block_size cache_size)
benchmarks=(blackscholes bodytrack fluidanimate)

for protocol in "${protocols[@]}"
do
  for dimension in "${dimensions[@]}"
  do
    for benchmark in "${benchmarks[@]}"
    do
      ./consolidate.sh ${protocol}/${dimension}/${benchmark} out ${protocol}_${dimension}_${benchmark}.csv
    done
  done
done
