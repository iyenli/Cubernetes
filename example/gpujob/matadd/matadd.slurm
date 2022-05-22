#!/bin/bash

# Do NOT edit the next 2 lines, or you will NOT get any output
#SBATCH --output=%j.out
#SBATCH --error=%j.err
###################################

#SBATCH --job-name=cublas
#SBATCH --partition=dgx2
#SBATCH -N 1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=6
#SBATCH --gres=gpu:1

ulimit -s unlimited
ulimit -l unlimited

module load gcc/8.3.0 cuda/10.1.243-gcc-8.3.0
nvcc matmult.cu -o matmult
./matmult