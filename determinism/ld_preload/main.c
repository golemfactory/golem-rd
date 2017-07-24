#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <time.h>
#include <assert.h>

#include "comdef.h"

int main()
{
    srand(42);
    int s = rand();
    assert(s == 71876166);

    uint32_t seed = (uint32_t) time(NULL);
    srand(seed);
    int r = rand();

    printf("random:  %d\n", r);
    printf("preseed: %d\n", s);

    if (r == s) {
        printf(COLOR_RED "YOU HAVE BEEN PWNED" COLOR_RESET "\n");
    }
    return 0;
}
