#ifndef __CHAIN_H__
#define __CHAIN_H__

#include <stddef.h>
#include <stdint.h>

typedef uint64_t capi_name;

typedef struct {
	uint8_t data[20];
} capi_checksum160;

typedef struct {
	uint8_t data[32];
} capi_checksum256;

typedef struct {
	uint8_t data[64];
} capi_checksum512;

#endif // __CHAIN_H__
