// Always compile with `gcc Zverala.c -lm'
#include<math.h>
#include<stdlib.h>
#include<stdio.h>
#include<string.h>

int this_doubleyear;
char direction;
int outward; //
int* digits;
int year_number_length;

int count_number_length(int number) {
    int count = 0;
    while (number != 0) {
        number /= 10;
        count++;
    }
    //printf("count: %d\n", count);
    return count;
}

void fill_in_digits() {
    int this_year = this_doubleyear;
    int digit_count = count_number_length(this_year);
    
    digits = (int *) malloc ((digit_count + 1) * sizeof(int));
    int* digits_temp = digits;
    year_number_length = digit_count;
    //printf("Pointless, but necessary for some fucked-up reason: %d\n", digit_count);
    this_year = this_doubleyear;
    printf("double year: %d %c\n", this_doubleyear, direction);
    for (int i = 0; i < digit_count; i++) {
        int num = this_year % 10;
        *(digits_temp + (digit_count - 1) - i) = this_year % 10;
        /* for (int j = 0; j <= i; j++) {
            printf("%d: %d\n", j, *(digits + j*sizeof(int)));
        } */
        //digits_temp += 1;//sizeof(int);
        this_year /= 10;
    }
}

int* convert_int(int number) {
    const int length = count_number_length(number);
    char string[length+1]; // +1 for \x00
    sprintf(string, "%d", number);
    int* digs = (int *) malloc (sizeof(int) * length);
    //int digs_ar[length];
    int* dig_temp = digs;
    for (int i = 0; i < length; i++) {
        //digs_ar[i] = string[i] - '0';
        //printf("Here%d: %d\n    %d\n", i, digs_ar[i], string[i] - '0');
        *dig_temp = string[i] - '0';
        dig_temp += 1;
    }
    //digs = digs_ar;
    return digs;
}

int* convert_double(double number, int decimal_points) {
    number = number * decimal_points;
    int num = (int) number;
    printf("truncated number: %d\n", num);
    return convert_int(num);
}

int calculate_a() {
    return this_doubleyear % 9 + 1;
}

int convert_int_arr_to_int(int* array, int array_length) {
    int result = 0;
    for (int i = array_length - 1; i >= 0; i--) {
        result += ((int) pow(10.0, i)) * array[array_length - i - 1];
        //printf("result: %d\n", result);
    }
    printf("int from array: %d\n", result);
    return result;
}

int calculate_b() {
    // Sub-calculation 1
    // (-1; +1)
    int first_sub_calc[year_number_length];
    for (int i = 0; i < year_number_length; i++) {
        if (i % 2 == 1) {
            int help = *(digits + i) + 1;
            if (help >= 10) {
                help = 1;
            }
            first_sub_calc[i] = help;
        } else {
            int help = *(digits + i) - 1;
            if (help <= 0) {
                help = 9;
            }
            first_sub_calc[i] = help;
        }
    }

    // Sub-calculation 2
    // (+1; -1)
    int second_sub_calc[year_number_length];
    for (int i = 0; i < year_number_length; i++) {
        if (i % 2 == 1) {
            int help = *(digits + i) - 1;
            if (help <= 0) {
                help = 9;
            }
            second_sub_calc[i] = help;
        } else {
            int help = *(digits + i) + 1;
            if (help >= 10) {
                help = 1;
            }
            second_sub_calc[i] = help;
        }
    }

    printf("(-1;+1); (+1; -1):\n");
    for (int i = 0; i < year_number_length; i++) {
        printf("first: [%d]; second: [%d]\n", first_sub_calc[i], second_sub_calc[i]);
    }

    int first_sub_num = convert_int_arr_to_int(first_sub_calc, year_number_length);
    int second_sub_num = convert_int_arr_to_int(second_sub_calc, year_number_length);

    // Sub-calculation 3
    double third_sub_num;
    if (outward) {
        third_sub_num = (double) first_sub_num / second_sub_num;
    } else {
        third_sub_num = (double) second_sub_num / first_sub_num;
    }
    printf("third: [%lf]\n", third_sub_num);

    int* sub_result = convert_double(third_sub_num, 4); // 4 is the decimal point count for years.
    // TODO Add digits together.
    return 0;
}

int calculate_c() {
    return 0;
}

int dragons_is_present() {
    return 1;
}

void convert_year() {
    const int beginning = -52131; // Normal year when this cycle's Klvanistic calendar started
    char normal_year[20]; // At most 4 digits, plus a potential negative, plus newline
    scanf("%s", normal_year);
    int this_year = atoi(normal_year); // Converts string to number, ignores all which is not a number
    int difference = this_year - beginning; // No. of normal years since the beginning
    int double_year = difference / 2;
    outward = 0;
    direction = 'S';
    if (difference % 2 == 1) { // If the difference is odd, change accordingly
        double_year ++;
        outward = 1;
        direction = 'O';
    }

    // Set global variables
    this_doubleyear = double_year;
}

void free_memory() {
    free(digits);
}

int main(int argc, char const *argv[])
{
    convert_year(); // Sets global variables.
    fill_in_digits(); // Split year into digits and store them globaly
    int a = calculate_a();
    int b = calculate_b();
    /* int c = calculate_c();
    bool dragons = dragons_is_present(); */
    
    printf("year_number_length: %d\n", year_number_length);
    for (int i = 0; i < year_number_length; i++) { // Test fill_in_digits
        printf("[%d]\n", *(digits + i));// * sizeof(int)));
    }
    printf("Next: test convert_int\n");

    int a_test = 72346; // Test convert_int
    int* ar = convert_int(a_test);
    for (int i = 0; i < 5; i++) {
        printf("[%d]\n", *(ar + i));//*sizeof(int)));
    }

    printf("Next: test convert_double\n");
    double b_test = 2.45639; // Test for convert_double
    int decimal = 10000;
    int* c_ar = convert_double(b_test, decimal);
    for (int i = 0; i < 5; i++) {
        printf("[%d]\n", *(c_ar + i));
    }

    free(c_ar);
    free(ar);
    free_memory();
    return 0;
}