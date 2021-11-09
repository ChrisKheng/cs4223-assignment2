#!/bin/bash
input_dir=$1
output_dir=$2
output_name=$3

if [[ $# -lt 3 ]]
then
  echo "Usage: consolidate [input dir] [output_dir] [output name]"
  exit 1
fi

mkdir -p $output_dir

stats=(${input_dir}/*.err)

for file in "${stats[@]}"
do
  filename=$(basename $file)
  filename="${filename%.*}"
  head -n1 $file | awk -F ',' -v name="${filename}" 'BEGIN{OFS=",";} {print name,$2;}' >> $output_name
  sort -t ',' -k1 -n $output_name -o $output_name
done

mv $output_name $output_dir
