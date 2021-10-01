#ifndef __DB_H__
#define __DB_H__

#include <stdint.h>

typedef struct {
	uint64_t lo;
	uint64_t hi;
} uint128_t;

//typedef long double long_double;

typedef struct {
	uint64_t lo;
	uint64_t hi;
} float128_t;




#endif // __DB_H__