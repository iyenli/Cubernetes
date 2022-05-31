#include <cuda_runtime.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#define N 8
#define IDX2C(i, j, ld) (((j) * (ld)) + (i))
#define BLOCK_SIZE 8

__global__ void matrixAdd(const float** A, const float** B, float** C,
    int M, int N)
{
    int i = blockIdx.x * blockDim.x + threadIdx.x;
    int j = blockIdx.y * blockDim.y + threadIdx.y;
    if (i < N && j < N)
        C[i][j] = A[i][j] + B[i][j];
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
    matrixAdd<<<dimGrid, dimBlock>>>(devPtrA, devPtrB, devPtrC, N, N);

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
