#include<stdio.h>
#include<stdlib.h>

int main(int argc, char const *argv[])
{
    int* digit = (int *) malloc (4 * sizeof(int));
    int a = 1234;
    int length = a % 10;

    int* digit_temp = digit;

    for (int i = length - 1; i >= 0; i--) {
        *digit_temp = a % 10;
        digit_temp += 1;
        a /= 10;
    }

    for (int i = 0; i < length; i++) {
        printf("%d\n", *(digit + i));
    }
    free(digit);
    return 0;
}
