#include <stdio.h>
#include <stdint.h>

typedef __int128 int128;

void int128_from_int64(int128* a, int64_t* b) {
    *a = (int128)(*(int64_t*)(b));
}

void int128_to_int64(int128* a, int64_t* b) {
    *b = (int64_t)(*(int128*)(b));
}

void int128_add(int128* a, int128* b, int128* c) {
    *c = *a + *b;
}

void int128_sub(int128* a, int128* b, int128* c) {
    *c = *a - *b;
}

// void int128_abs(int128* a, int128* b) {
//     if (*a > 0) {
//         *b = *a;
//     } else {
//         *b = -*a;
//     }
// }

void int128_mul(int128* a, int128* b, int128* c) {
    *c = *a * *b;
}

void int128_div(int128* a, int128* b, int128* c) {
    *c = *a / *b;
}

int int128_cmp(int128* a, int128* b) {
    if (*a > *b) {
        return 1;
    } else if (*a < *b) {
        return -1;
    } else {
        return 0;
    }
}
