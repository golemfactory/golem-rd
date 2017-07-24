#include <time.h>
#include <stdio.h>
#include "comdef.h"

time_t time(time_t* tloc)
{
    (void) tloc;

    printf("time: " COLOR_RED "YOU HAVE BEEN PWNED!!\n" COLOR_RESET);
    return 42;
}
