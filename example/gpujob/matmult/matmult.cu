#include <cuda_runtime.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#define N 8
#define IDX2C(i, j, ld) (((j) * (ld)) + (i))
#define BLOCK_SIZE 8

__global__ void matrixMulSharedKernel_op1(float* fpMatrixA, float* fpMatrixB,
    float* fpMatrixC, int m, int n, int k)
{
    int nRow = blockIdx.y * blockDim.y + threadIdx.y;
    int nCol = blockIdx.x * blockDim.x + threadIdx.x;
    float fCVal = 0.0f;

    __shared__ float shTileA[BLOCK_SIZE][BLOCK_SIZE];
    __shared__ float shTileB[BLOCK_SIZE][BLOCK_SIZE];

    int nIter = (k + BLOCK_SIZE - 1) / BLOCK_SIZE;
    for (int i = 0; i < nIter; i++) {
        // load data from global memory to shared memory
        shTileA[threadIdx.y][threadIdx.x] = fpMatrixA[nRow * k + i * BLOCK_SIZE + threadIdx.x];
        shTileB[threadIdx.y][threadIdx.x] = fpMatrixB[(i * BLOCK_SIZE + threadIdx.y) * n + nCol];

        // sync to wait for all threads in one block to finish loading datas
        __syncthreads();

        // sub-matrix multiply
        for (int l = 0; l < BLOCK_SIZE; l++) {
            fCVal += shTileA[threadIdx.y][l] * shTileB[l][threadIdx.x];
        }

        // sync to wait for all threads in one block to finish compute
        __syncthreads();
    }

    // store results into global memory
    fpMatrixC[nRow * n + nCol] = fCVal;
}

void printfMatrix(float* a, int m, int n)
{
    for (int j = 0; j < m; j++) {
        printf("[");
        for (int i = 0; i < n; i++) {
            printf("\t%lg", a[IDX2C(i, j, m)]);
        }
        printf("\t]\n");
    }
}

int main(void)
{
    int i, j;
    float* devPtrA;
    float* devPtrB;
    float* devPtrC;
    float* a = 0;
    float* b = 0;
    float* c = 0;
    a = (float*)malloc(N * N * sizeof(*a));
    b = (float*)malloc(N * N * sizeof(*b));
    c = (float*)malloc(N * N * sizeof(*b));

    for (j = 0; j < N; j++) {
        for (i = 0; i < N; i++) {
            a[IDX2C(i, j, N)] = (float)(i * N + j + 1);
            b[IDX2C(i, j, N)] = (float)(j * N + i + 1);
        }
    }
    cudaMalloc((void**)&devPtrA, N * N * sizeof(*a));
    cudaMalloc((void**)&devPtrB, N * N * sizeof(*b));
    cudaMalloc((void**)&devPtrC, N * N * sizeof(*c));

    cudaMemcpy(devPtrA, a, sizeof(*a) * N * N, cudaMemcpyHostToDevice);
    cudaMemcpy(devPtrB, b, sizeof(*b) * N * N, cudaMemcpyHostToDevice);

    dim3 dimBlock(BLOCK_SIZE, BLOCK_SIZE);
    dim3 dimGrid((int)ceil(N / BLOCK_SIZE), (int)ceil(N / BLOCK_SIZE));
    matrixMulSharedKernel_op1<<<dimGrid, dimBlock>>>(devPtrA, devPtrB, devPtrC, N, N, N);

    cudaMemcpy(c, devPtrC, sizeof(*c) * N * N, cudaMemcpyDeviceToHost);

    printf("A x B = C\n");
    printf("\nA:\n");
    printfMatrix(a, N, N);
    printf("\nB:\n");
    printfMatrix(b, N, N);
    printf("\nC:\n");
    printfMatrix(c, N, N);

    cudaFree(devPtrA);
    cudaFree(devPtrB);
    cudaFree(devPtrC);

    free(a);
    free(b);
    free(c);
    return EXIT_SUCCESS;
}
