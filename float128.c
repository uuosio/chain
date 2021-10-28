#include <stdio.h>

typedef long double float128;

void float128_from_double(float128* a, double* b) {
    *a = (float128)(*(double*)(b));
}

void float128_to_double(float128* a, double* b) {
    *b = (double)(*(float128*)(b));
}

void float128_add(float128* a, float128* b, float128* c) {
    *c = *a + *b;
}

void float128_sub(float128* a, float128* b, float128* c) {
    *c = *a - *b;
}

void float128_abs(float128* a, float128* b) {
    if (*a > 0) {
        *b = *a;
    } else {
        *b = -*a;
    }
}

void float128_mul(float128* a, float128* b, float128* c) {
    *c = *a * *b;
}

void float128_div(float128* a, float128* b, float128* c) {
    *c = *a / *b;
}

int float128_cmp(float128* a, float128* b) {
    if (*a > *b) {
        return 1;
    } else if (*a < *b) {
        return -1;
    } else {
        return 0;
    }
}
