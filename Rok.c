#include<stdlib.h>
#include<stdio.h>

int main(int argc, char const *argv[])
{
    const int beginning = -52131; // Normal year when this cycle's Klvanistic calendar started
    char normal_year[6]; // At most 4 digits, plus a potential negative, plus newline
    scanf("%s", normal_year);
    int this_year = atoi(normal_year); // Converts string to number, ignores all which is not a number
    int difference = this_year - beginning; // No. of normal years since the beginning
    int double_year = difference / 2;
    char outward = 'S';
    if (difference % 2 == 1) { // If the difference is odd, change accordingly
        double_year ++;
        outward = 'O';
    }

    printf("%d %c\n", double_year, outward);
    return 0;
}
