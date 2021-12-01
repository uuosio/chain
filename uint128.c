#include <stdio.h>
#include <stdint.h>

typedef unsigned __int128 uint128;

void uint128_from_uint64(uint128* a, uint64_t* b) {
    *a = (uint128)(*(uint64_t*)(b));
}

void uint128_to_uint64(uint128* a, uint64_t* b) {
    *b = (uint64_t)(*(uint128*)(b));
}

void uint128_add(uint128* a, uint128* b, uint128* c) {
    *c = *a + *b;
}

void uint128_sub(uint128* a, uint128* b, uint128* c) {
    *c = *a - *b;
}

// void uint128_abs(uint128* a, uint128* b) {
//     if (*a > 0) {
//         *b = *a;
//     } else {
//         *b = -*a;
//     }
// }

void uint128_mul(uint128* a, uint128* b, uint128* c) {
    *c = *a * *b;
}

void uint128_div(uint128* a, uint128* b, uint128* c) {
    *c = *a / *b;
}

int uint128_cmp(uint128* a, uint128* b) {
    if (*a > *b) {
        return 1;
    } else if (*a < *b) {
        return -1;
    } else {
        return 0;
    }
}
