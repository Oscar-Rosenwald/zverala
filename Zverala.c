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
int *sins;
int *cosins;
long double *help_trigonometry;
const int ROOT_SIZE = 30;
int *intersections;

typedef enum {
    January,
    February,
    March,
    April,
    May,
    June,
    July,
    August,
    September,
    October,
    November,
    December
} Month;

struct Date {
    int day;
    Month month;
    int year;
    int doubleyear;
    int outward;
    char* dyear;
};

int days_of_month(Month month, int year) {
    switch (month)
    {
    case January:
    case March:
    case May:
    case July:
    case August:
    case October:
    case December:
        return 31;
    case April:
    case June:
    case September:
    case November:
        return 30;
    default: // February
        if (year % 4 == 0) {
            if (year % 100 == 0 && year % 1000 != 0) {
                return 28;
            } else if (year % 1000 == 0) {
                return 29;
            } else {
                return 28;
            }
        } else {
            return 28;
        }
    }
}

Month next_month(Month month) {
    switch (month)
    {
    case January:
        return February;
    case February:
        return March;
    case March:
        return April;
    case April:
        return May;
    case May:
        return June;
    case June:
        return July;
    case July:
        return August;
    case August:
        return September;
    case September:
        return October;
    case October:
        return November;
    case November:
        return December;    
    default:
        return January;
    }
}

void add_days_to_date(struct Date *date, int days) {
    while (days != 0) {
        Month month = date->month;
        int day = date->day;
        int days_in_month = days_of_month(month, date->year);

        int difference = days_in_month - day; // No. of days remaining in month
        if (difference >= days) {
            date->day = day + days; // Only add the number of days to add, but stay in the same month
            days = 0; // This will stop the loop
        } else {
            days = days - difference - 1; // -1 To go to another month
            date->day = 1;
            date->month = next_month(month);
            if (date->month == January) {
                date->year++;
                if (date->outward) {
                    date->outward = 0;
                    sprintf(date->dyear, "%dS", date->doubleyear);
                } else {
                    date->doubleyear++;
                    date->outward = 1;
                    sprintf(date->dyear, "%dO", date->doubleyear);
                }
            }
        }
    }
}

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
    printf("Fill_in_digits - double year: %d %c\n", this_doubleyear, direction);
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
    number = number * (pow(10, decimal_points));
    int num = (int) number;
    printf("truncated number: %d\n", num);
    return convert_int(num);
}

int calculate_a() {
    int result = this_doubleyear % 9 + 1;
    printf("a: %d\n", result);
    return result;
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

int add_array_integers(int* array, int array_length) {
    int result = 0;
    for (int i = 0; i < array_length; i++) {
        printf("array to add: %d\n", array[i]);
        result += array[i];
    }
    printf("adding array: %d\n", result);
    if (result >= 10) {
        int* new_operation = convert_int(result);
        int length = count_number_length(result);
        result = add_array_integers(new_operation, length);
        free(new_operation);
    }
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
    printf("third: %lf\n", third_sub_num);

    int number_after_decimal = year_number_length;
    int sub_result_added;

    int* sub_result = convert_double(third_sub_num, number_after_decimal); // 4 is the decimal point count for years.
    // Check the length of the number
    if (third_sub_num < 1) { // Length of number -> 4
        sub_result_added = add_array_integers(sub_result, number_after_decimal);
    } else { // Length of number -> 5
        sub_result_added = add_array_integers(sub_result, number_after_decimal + 1);
    } // Never exeeds five.
    printf("b: %d\n", sub_result_added);
    free(sub_result);
    return sub_result_added;
}

void swap_numbers(int* one, int* two) {
    printf("swapping: %d, %d\n", *one, *two);
    int help = *one;
    *one = *two;
    *two = help;
}

int calculate_c() {
    int original_places[year_number_length]; // [2][1][4][3] (places)
    int current_min;
    int previous_min = 0; // Start at 0 so the first iteration passes.
    int prev_min_index = -1; // Start at -1 so the first iteration always passes.
    int prev_min_index_help = prev_min_index;
    printf("\nNext: c test\n");
    for (int places = 0; places < year_number_length; places++) {
        current_min = 10; // We start with this, so that any digits will be smaller on the first iteration of digs.
        for (int digs = 0; digs < year_number_length; digs++) {
            int* num = digits + digs;
            //printf("*num = %d\ncurrent_min = %d\nprevious_min = %d\nprev_min_index = %d\n\n", *num, current_min, previous_min, prev_min_index);
            if (*num < current_min) {
                if (*num > previous_min) {
                    current_min = *num;
                    original_places[places] = digs + 1;
                    prev_min_index_help = digs;
                } else if (*num == previous_min) {
                    if (digs > prev_min_index) {
                        current_min = *num;
                        original_places[places] = digs + 1;
                        prev_min_index_help = digs;
                    }
                }
            }
        }
        prev_min_index = prev_min_index_help;
        previous_min = current_min;
    }

    for (int i = 0; i < year_number_length; i++) {
        printf("[%d] ", original_places[i]);
    }

    int sub_result1 = convert_int_arr_to_int(original_places, year_number_length);
    int sub_result2 = this_doubleyear * sub_result1;
    return sub_result2 % 9 + 1;
}

int there_be_dragons() {
    int last_digit = this_doubleyear % 10;
    int dividers[] = {3, 5, 7};
    if (outward) {
        if (last_digit == 7) {
            dividers[1] = 0;
        } else {
            dividers[1] = last_digit;
        }
    } else {
        if (last_digit == 5) {
            dividers[2] = 0;
        } else {
            dividers[2] = last_digit;
        }
    }

    int how_many_dividers = 0;

    for (int i = 0; i < sizeof(dividers) / sizeof(dividers[0]); i++) {
        if (dividers[i] != 0) {
            if (this_doubleyear % dividers[i] == 0) {
                printf("%d / %d = %d\n", this_doubleyear, dividers[i], this_doubleyear / dividers[i]);
                how_many_dividers++;
            }
        }
    }

    if (how_many_dividers == 2) {
        return 1;
    }
    return 0;
}

int convert_year(int global, int this_year) {
    const int beginning = -52131; // Normal year when this cycle's Klvanistic calendar started
    int difference = this_year - beginning; // No. of normal years since the beginning
    int double_year = difference / 2;
    if (global) {
        printf("global\n");
        outward = 0;
        direction = 'S';
    }
   if (difference % 2 == 1) { // If the difference is odd, change accordingly
        double_year ++;
        if (global) {
            outward = 1;
            direction = 'O';
        }
    }
    if (global) {
        this_doubleyear = double_year;
    }
    return double_year;
}

int days_of_kyear(int year, int first_solstice, int second_solstice) {
    struct Date *date = malloc(sizeof(struct Date));
    date->year = year;
    date->doubleyear = convert_year(0, year);
    char outward = 'O';
    if (year % 2 == 0) {
        date->outward = 1;
    } else {
        date->outward = 0;
        outward = 'S';
    }
    date->day = first_solstice;
    date->month = December;
    char* dyear = malloc(16*sizeof(char));
    sprintf(dyear, "%d%c", date->doubleyear, outward);
    date->dyear = dyear;

    int result = 350; // Start with a safe but sufficiently large value

    add_days_to_date(date, result);
    while ((date->day != second_solstice) || (date->day == second_solstice && date->month != December)) {
        add_days_to_date(date, 1);
        result++;
    }

    printf("doubleyear for solstices: %s", date->dyear);

    free(dyear);
    free(date);
    return result;
}

void convert_year_globally() {
    char normal_year[20]; // At most 4 digits, plus a potential negative, plus newline
    scanf("%s", normal_year);
    int this_year = atoi(normal_year); // Converts string to number, ignores all which is not a number
    printf("normal year: %s - %d\n", normal_year, this_year);
    convert_year(1, this_year);
}

void free_memory() {
    free(digits);
    free(sins);
    free(cosins);
    free(help_trigonometry);
    free(intersections);
}

void get_sins(int a) {
    for (int k = 0; k < ROOT_SIZE; k++) {
        *(help_trigonometry + k - 1) = (long double) sqrt((M_PI *  k)/a);
    }
}

long double calculate_period (int k, int b) {
    return (M_PI * k) / (2*b);
}

void get_cosins(int b, int c) {
    long double root;
    int parameter_bound = -1;
    do {
        parameter_bound += 2;
        root = calculate_period(parameter_bound, b);
        printf("root = %Lf; k = %d\n", root, parameter_bound);
    } while (root < 1);
    parameter_bound -= 2; // After the last check, parameter_bound will be 2 too big.

    int how_many_roots = (parameter_bound / 2 + 1);
    long double sub_results[how_many_roots][2];
    // OBSOLITE 2 for k (one for PI - result); one for -k (not two - x would be negative, unless PI - (- k))

    for (int i = 0; i < how_many_roots; i++) {
        int parameter = i * 2 + 1; // Get the odd number
        printf("period_max = %d; i = %d\nparameter = %d\n", parameter_bound, i, parameter);
        root = calculate_period(parameter, b);
        long double argsin = asinl(root);
        long double res1 = argsin / c;
        long double res2 = (M_PI - argsin) / c;
        //root = calculate_period(-parameter, b);
        //argsin = asinl(root);
        //long double res3 = (M_PI - argsin) / c;
        sub_results[i][0] = res1;
        sub_results[i][1] = res2;
        //sub_results[i][2] = res3;
        printf("cosin: [%Lf]\ncosin: [%Lf]\n", sub_results[i][0], sub_results[i][1]);
    }

    printf("ROOT_SIZE / (2 * how_many_roots) = %d\n", ROOT_SIZE / (2 * how_many_roots));
    
    int repeat = ROOT_SIZE / (2 * how_many_roots);
    for (int period = 0; period < repeat; period++) {
        for (int j = 0; j < how_many_roots; j++) {
            int offset = period * how_many_roots * 2;
            *(help_trigonometry + offset + j) = sub_results[j][0] + (period * (M_PI / c));
            *(help_trigonometry + offset + how_many_roots * 2 - j - 1) = sub_results[j][1] + (period * (M_PI / c));
            printf("%d  b: [%Lf]\n   b: [%Lf]\n", period, help_trigonometry[offset + j], help_trigonometry[offset + how_many_roots * 2 - j - 1]);
        }
    }
}

void expand_double_to_int_array (int *array_to, long double *array_from, int array_length) {
    for (int i = 0; i < array_length; i++) {
        int decimal;
        if (*(array_from + i) >= 1) {
            decimal = 5;
        } else if (array_from[i] < 0.1) {
            decimal = 3;
        } else {
            decimal = 4;
        }
        int * converted = convert_double(*(array_from + i), 4);
        *(array_to + i) = convert_int_arr_to_int(converted, decimal);
        free(converted);
    }
}

int main(int argc, char const *argv[])
{
    convert_year_globally(); // Sets global variables.
    printf("Global - this doubleyear: %d\n", this_doubleyear);
    fill_in_digits(); // Split year into digits and store them globaly
    
    int a = calculate_a();
    int b = calculate_b();
    int c = calculate_c();
    int dragons = there_be_dragons();

    sins = malloc(ROOT_SIZE*sizeof(long double));
    cosins = malloc(ROOT_SIZE*sizeof(long double));
    help_trigonometry = malloc(ROOT_SIZE*sizeof(long double));

    get_sins(a);
    expand_double_to_int_array(sins, help_trigonometry, ROOT_SIZE);
    for (int i = 0; i < ROOT_SIZE; i++) {
        printf("sin %d: %d\n", i, *(sins + i));
    }

    get_cosins(b, c);
    expand_double_to_int_array(cosins, help_trigonometry, ROOT_SIZE);
    for (int i = 0; i < ROOT_SIZE; i++) {
        printf("sin %d: %d\n", i, *(sins + i));
    }
    for (int i = 0; i < ROOT_SIZE; i++) {
        printf("cosin %d: %d\n", i, *(cosins + i));
    }

    int intersection_count = 34;
    intersections = malloc(intersection_count*sizeof(int));
    *intersections = 0;
    // Get 34 results - 1 for 0;
    int sin = 0;
    int cos = 0;
    for (int i = 1; i < intersection_count; i++) {
        if (*(sins + sin) > *(cosins + cos)) {
            *(intersections + i) = *(cosins + cos);
            cos++;
        } else {
            *(intersections + i) = *(sins + sin);
            sin++;
        }
    }
    for (int i = 0; i < intersection_count; i++) {
        printf("intersecion %d: %d\n", i, *(intersections + i));
    }

    if (dragons) {
        printf("There be dragons!\n");
    } else {
        printf("No dragons this year.\n");
    }

    printf("c = %d\n", c);

    
    
    
    
    //-----------------------------------------------------------------------------------------------------
    //-----------------------------------------------------------------------------------------------------
    //-----------------------------------------------------------------------------------------------------
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
    int decimal = 4;
    int* c_ar = convert_double(b_test, decimal);
    for (int i = 0; i < 5; i++) {
        printf("[%d]\n", *(c_ar + i));
    }

    printf("\nNext: add days test:\n");
    struct Date *date = malloc(sizeof(struct Date));
    date->year = 2013;
    date->doubleyear = 27072;
    date->outward = 0;
    char* double_year_test = malloc(16*sizeof(char));
    sprintf(double_year_test, "27072S");
    date->dyear = double_year_test;
    date->month = March;
    date->day = 21;

    add_days_to_date(date, 378+370);
    printf("febuary: %d\n", days_of_month(February, 2014));
    printf("year: %d\ndoubleyear: %d%d\ndyear: %s\nmonth: %d; day: %d\n", date->year, date->doubleyear, date->outward, date->dyear,
    date->month, date->day);

    printf("\nNext: test of kyear\ndays in kyear: %d\n", days_of_kyear(2020, 21, 23));

    free(c_ar);
    free(ar);
    free(date);
    free(double_year_test);
    free_memory();
    return 0;
}