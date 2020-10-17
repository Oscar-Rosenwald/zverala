#include<stdlib.h>
#include<stdio.h>

int main(int argc, char const *argv[])
{
    int a = 1234;
    int length = a % 10;
    //printf("a: %d, length: %d\n", a, length);

    int* digit = malloc(sizeof(int)*length);
    printf("memory: %p\n", digit + sizeof(int)*length);
    for (int i = 1; i < length + 1; i++) {
        *(digit + (i - 1) * sizeof(int)) = i;
        //printf("i: %d\n", i);
        //printf("digit: %d\n\n", *(digit + (i - 1) * sizeof(int)));
    }

    for (int i = 0; i < length; i++) {
        int num = a % 10;
        *(digit + i) = num;
        //printf("pointer: %p\ndigit: %d\n", digit + i * sizeof(int), *(digit + i * sizeof(int)));
    //    a = a / 10;
        //printf("a: %d\n\n", a);
    }

    printf("size of int: %ld\n", sizeof(int));

    for (int i = 0; i < length; i++) {
        int* poin = digit + i;
        printf("pointer: %p\n%d\n", poin, *(poin));
    }

    //printf("3: %d\n", *(digit + 2 * sizeof(int)));
    //printf("3: %d\n", *(digit + 2 * sizeof(int)));
    //printf("3: %d\n", *(digit + 2));

    return 0;
}
