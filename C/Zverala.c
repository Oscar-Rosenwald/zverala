// Always compile with `gcc Zverala.c -lm'
#include<math.h>
#include<stdlib.h>
#include<stdio.h>
#include<string.h>

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

int this_doubleyear;
int this_year;
char direction;
int outward; //
int* digits;
int year_number_length;
int *sins;
int *cosins;
long double *help_trigonometry;
char file_name[100];
const int ROOT_SIZE = 36;
int *intersections;
struct Being *mythical_beings;
char *beings[] = {
    "Chiméra",

    "Fénix",
    "Toch Amogaši",
    
    "Sfinga",
    "Vlkodlak",
    "Zlovlk",
    
    "Jednorožec",
	"Griffin",
    "Lví želva",

    "Kraken",
    "Kyklop",
    "Syréna",
	"Yeti",
    "Nessie",
    "Vyjící chluporyba",

    "Olifant",
    "Ždiboň",
    "Ent",
    "Labuť",

    "Kerberos",
    "Bazilišek",
    "Akromantule",
    "Goa'uld",
    "Vetřelec",

    "Létající bizon",
    "Pegas",
    "Mothra",
    "Sleipnir",
    "Velká A'tuin",
    "Horus",

    "Cthulhu",
    "Hydra",
    "Balrog",
    "Odgru Jahad",
};
char *dragons_types[] = {
    "ohně",
    "země",
    "života",
    "vody",
    "dřeva",
    "smrti",
    "vzduchu",
    "chaosu"
};
int dragon_after_index[] = {0, 2, 5, 8, 14, 18, 23, 29};

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
            if (year % 1000 == 0) {
                return 29;
            } else if (year % 100 == 0) {
                return 28;
            } else {
                return 29;
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
    return count;
}

void fill_in_digits() {
    int this_year = this_doubleyear;
    int digit_count = count_number_length(this_year);
    
    digits = (int *) malloc ((digit_count + 1) * sizeof(int));
    int* digits_temp = digits;
    year_number_length = digit_count;
    this_year = this_doubleyear;
    for (int i = 0; i < digit_count; i++) {
        int num = this_year % 10;
        *(digits_temp + (digit_count - 1) - i) = this_year % 10;
        this_year /= 10;
    }
}

int* convert_int(int number) {
    const int length = count_number_length(number);
    char string[length+1]; // +1 for \x00
    sprintf(string, "%d", number);
    int* digs = (int *) malloc (sizeof(int) * length);
    int* dig_temp = digs;
    for (int i = 0; i < length; i++) {
        *dig_temp = string[i] - '0';
        dig_temp += 1;
    }
    return digs;
}

int* convert_double(double number, int decimal_points) {
    number = number * (pow(10, decimal_points));
    int num = (int) number;
    return convert_int(num);
}

int calculate_a() {
    int result = this_doubleyear % 9 + 1;
    return result;
}

int convert_int_arr_to_int(int* array, int array_length) {
    int result = 0;
    for (int i = array_length - 1; i >= 0; i--) {
        result += ((int) pow(10.0, i)) * *(array + array_length - i - 1);
    }
    return result;
}

int add_array_integers(int* array, int array_length) {
    int result = 0;
    for (int i = 0; i < array_length; i++) {
        result += array[i];
    }
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

    int first_sub_num = convert_int_arr_to_int(first_sub_calc, year_number_length);
    int second_sub_num = convert_int_arr_to_int(second_sub_calc, year_number_length);

    // Sub-calculation 3
    double third_sub_num;
    if (outward) {
        third_sub_num = (double) first_sub_num / second_sub_num;
    } else {
        third_sub_num = (double) second_sub_num / first_sub_num;
    }

    int number_after_decimal = year_number_length;
    int sub_result_added;

    int* sub_result = convert_double(third_sub_num, number_after_decimal); // Creates an array with digits up until the number_after_decimal-th digit after the decimal point.
    // Check the length of the number
    if (third_sub_num < 1) { // Only the digits after decimal point
        sub_result_added = add_array_integers(sub_result, number_after_decimal);
    } else { // Use the digit before the decimal
        sub_result_added = add_array_integers(sub_result, number_after_decimal + 1);
    } // Never more than one digit before decimal.
    free(sub_result);
    /* printf("sub_result_added = %d\n", sub_result_added); */
    return sub_result_added;
}

void swap_numbers(int* one, int* two) {
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
    for (int places = 0; places < year_number_length; places++) {
        current_min = 10; // We start with this, so that any digits will be smaller on the first iteration of digs.
        for (int digs = 0; digs < year_number_length; digs++) {
            int* num = digits + digs;
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
    sprintf(dyear, "%d %c", date->doubleyear, outward);
    date->dyear = dyear;

    int result = 350; // Start with a safe but sufficiently large value

    add_days_to_date(date, result);
    while ((date->day != second_solstice) || (date->day == second_solstice && date->month != December)) {
        add_days_to_date(date, 1);
        result++;
    }

    free(dyear);
    free(date);
    return result;
}

void printYearFromFile (FILE* file) {
    char temp[100];

    fgets(temp, 100, file); // Empty line
    printf("Rok %d nalezen v souboru %s\n", this_year, file_name);
    while(fgets(temp, 100, file)) {
        if (strncmp(temp, "Rok", 3) == 0) {
            return;
        }
        printf("%s", temp);
    }
}

void convert_year_globally() {
    char normal_year[4]; // At most 4 digits, plus a potential negative, plus newline
	printf("Zadejte rok (normalni): ");
    scanf("%s", normal_year);
    char n_year[9];
    sprintf(n_year, "Rok %s", normal_year);
    this_year = atoi(normal_year); // Converts string to number, ignores all which is not a number
    FILE* file;
    char temp[100];
    if ((file = fopen(file_name, "r"))) {
        while (fgets(temp, 100, file)) {
            if (strncmp(temp, n_year, 8) == 0) {
                printYearFromFile(file);
                exit(0);
            }
        }
    }
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
    for (int k = 1; k < ROOT_SIZE; k++) {
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
    } while (root < 1);
    parameter_bound -= 2; // After the last check, parameter_bound will be 2 too big.

    int how_many_roots = (parameter_bound / 2 + 1);
    long double sub_results[how_many_roots][2];

    for (int i = 0; i < how_many_roots; i++) {
        int parameter = i * 2 + 1; // Get the odd number
        root = calculate_period(parameter, b);
        long double argsin = asinl(root);
        long double res1 = argsin / c;
        long double res2 = (M_PI - argsin) / c;
        sub_results[i][0] = res1;
        sub_results[i][1] = res2;
    }

    if (parameter_bound < 1) {
        for (int i = 0; i < ROOT_SIZE; i++) {
            *(help_trigonometry + i) = 1000;
        }
        return;
    }

    int repeat = ROOT_SIZE / (2 * how_many_roots);
    for (int period = 0; period < repeat; period++) {
        for (int j = 0; j < how_many_roots; j++) {
            int offset = period * how_many_roots * 2;
            *(help_trigonometry + offset + j) = sub_results[j][0] + (period * (M_PI / c));
            *(help_trigonometry + offset + how_many_roots * 2 - j - 1) = sub_results[j][1] + (period * (M_PI / c));
        }
    }
}

void expand_double_to_int_array (int *array_to, long double *array_from, int array_length) {
    for (int i = 0; i < array_length; i++) {
        int decimal;
        if (array_from[i] > 100) {
            decimal = 7;
        } else if (array_from[i] > 10) {
            decimal = 6;
        } else if (*(array_from + i) >= 1) {
            decimal = 5;
        } else if (array_from[i] < 0.1) {
            decimal = 3;
        } else {
            decimal = 4;
        }
        int * converted;
        if (decimal == 7) {
            converted = convert_double(*(array_from + i), 0);
        } else {
             converted = convert_double(*(array_from + i), 4);
        }
        *(array_to + i) = convert_int_arr_to_int(converted, decimal);
        free(converted);
    }
}

char* date_to_string(struct Date *date) {
    char *result = malloc(32 * sizeof (char));
    sprintf(result, "%d.%d. %d / %s", date->day, date->month + 1, date->year, date->dyear);
    return result;
}

void help() {
    printf("Jak na to:\n");
    printf("-n        ... Nepsat do souboru\n");
    printf("-f <xxx>  ... Hledat/zapsat do souboru xxx\n");
    printf("-h --help ... Zobrazit tento text\n");
}

int main(int argc, char const *argv[])
{
    sprintf(file_name, "Zverala.txt"); 
    if (argc > 1) {
        if (strcmp(argv[1], "-h") == 0 || strcmp(argv[1], "--help") == 0) {
            help();
            exit(1);
        }

        if (strcmp(argv[1], "-f") == 0) {
            if (argc == 3) {
                sprintf(file_name, "%s", argv[2]);
            } else {
                help();
                exit(1);
            }
        }
    }
    convert_year_globally(); // Sets changable global variables.
    fill_in_digits(); // Split year into digits and store them globaly
    
    int a = calculate_a();
    int b = calculate_b();
    int c = calculate_c();
    int dragons = there_be_dragons();
    printf("a: %d\nb: %d\nc: %d\n\n", a, b, c);

    sins = malloc(ROOT_SIZE*sizeof(long double));
    cosins = malloc(ROOT_SIZE*sizeof(long double));
    help_trigonometry = malloc(ROOT_SIZE*sizeof(long double));

    get_sins(a);
    expand_double_to_int_array(sins, help_trigonometry, ROOT_SIZE);

    get_cosins(b, c);
    expand_double_to_int_array(cosins, help_trigonometry, ROOT_SIZE);

    int intersection_count = 34;
    intersections = malloc(intersection_count*sizeof(int));
    *intersections = 0;
    // Get 34 results - 1 for 0;
    int sin = 0;
    int cos = 0;
    for (int i = 1; i < intersection_count; i++) {
        if (*(sins + sin) < *(cosins + cos)) {
            *(intersections + i) = *(sins + sin);
            sin++;
        } else {
            *(intersections + i) = *(cosins + cos);
            cos++;
        }
    }

    if (dragons) {
        printf("JE ROK DRAKU!\n");
    } else {
        printf("Letos bez draku.\n");
    }

    for (int i = 0; i < intersection_count - 1; i++) {
        *(intersections + i) = *(intersections + i + 1) - *(intersections + i);
    }

    int portions_added = 0;
    for (int i = 0; i < intersection_count - 1; i++) {
        portions_added += *(intersections + i);
    }

    int solstice1;
    int solstice2;
    printf("Zimni slunovrat roku %d (den v prosinci): ", this_year - 1);
    scanf("%d", &solstice1);
    printf("Zimni slunovrat roku %d (den v prosinci): ", this_year);
    scanf("%d", &solstice2);
    int kyear_lenth = days_of_kyear(this_year - 1, solstice1, solstice2);
    if (dragons) {
        kyear_lenth -= 8;
    }
    
    int used_days = 0;
    int *being_duration = malloc(34*sizeof(int));

	// Due to Chimera not always being first in the list,
	// we must set where the iterations over durations will go to
	// and where they'll stop. Then we'll assign the unused days
	// to Chimera, wherever it may be.
	int intersection_begin_index;
	int intersection_end_index;
	if (outward) { // Chimera is first this kyear
	  intersection_begin_index = 1;
	  intersection_end_index = intersection_count;
	} else { // Chimera is last
	  intersection_begin_index = 0;
	  intersection_end_index = intersection_count - 1;
	}
	  
    for (int i = intersection_begin_index; i < intersection_end_index; i++) {
        *(being_duration + i) = (kyear_lenth * *(intersections + i - intersection_begin_index)) / portions_added;
        used_days += *(being_duration + i);
    }

	if (outward) { // Assign unused days to Chimera in first place
	  *being_duration = kyear_lenth - used_days;
	} else { // Assign unused days to Chimera in last place
	  *(being_duration + intersection_count - 1) = kyear_lenth - used_days;
	}

    // In case we need to see the durations:
    /* for (int i = 0; i < intersection_count; i++) {
        printf("Being %d: %d days; address: %p\n", i, *(being_duration + i), being_duration + i);
    } */

    int dragons_index;
    if (outward) {
        dragons_index = 0;
    } else {
        dragons_index = 7;
    }
    struct Date *start_date = malloc(sizeof(struct Date));
    start_date->doubleyear = this_doubleyear;
    start_date->year = this_year - 1;
    start_date->outward = outward;
    start_date->month = December;
    start_date->day = solstice1;
    start_date->dyear = malloc(16*sizeof(char));
    sprintf(start_date->dyear, "%d %c", this_doubleyear, direction);

    // Prepare to write to file, if that is what's wanted
    int to_file = 1;
    if (argc > 1){
        to_file = strcmp(argv[1], "-n");
        if (to_file == 0) {
            to_file = 0;
        } else {
            to_file = 1;
        }
    }

    FILE *file;
    if (to_file) {
        file = fopen(file_name, "a");
        fprintf(file, "Rok %d / %d %c\n\n", this_year, this_doubleyear, direction);
    }

    for (int i = 0; i < intersection_count; i++) {
        if (*(being_duration + i) != 0) {
            char* date_string = date_to_string(start_date);

            if (outward) {
                if (to_file) {
                    fprintf(file, "%22s______%s\n", date_string, beings[i]);
                }
                printf("%22s______%s\n", date_string, beings[i]);
            } else {
                if (to_file) {
                    fprintf(file, "%22s______%s\n", date_string, beings[intersection_count - i - 1]);
                }
                printf("%22s______%s\n", date_string, beings[intersection_count - i - 1]);
            }

            add_days_to_date(start_date, *(being_duration + i));
            free(date_string);
        }
        if (dragons) {
            // If it's a turn for a dragon
            char* date_string = date_to_string(start_date);

            if (outward) {
                if (dragon_after_index[dragons_index] == i) {
                    if (to_file) {
                        fprintf(file, "%22s______Drak %s\n", date_string, dragons_types[dragons_index]);
                    }
                    printf("%22s______Drak %s\n", date_string, dragons_types[dragons_index]);
                    dragons_index++;
                    add_days_to_date(start_date, 1);
                }
            } else {
                if (dragon_after_index[dragons_index] == intersection_count - i - 2) {
                    if (to_file) {
                        fprintf(file, "%22s______Drak %s\n", date_string, dragons_types[dragons_index]);
                    }
                    printf("%22s______Drak %s\n", date_string, dragons_types[dragons_index]);
                    dragons_index--;
                    add_days_to_date(start_date, 1);
                }
            }
            free(date_string);
        }
    }
    if (to_file) {
        fclose(file);
    }
    free(start_date->dyear);
    free(start_date);
    free(being_duration);
    free_memory();
    return 0;
}