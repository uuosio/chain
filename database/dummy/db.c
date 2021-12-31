#include <stddef.h>
#include <stdint.h>
#include <stdio.h>

static void go_panic(char* funcName) {
    printf("%s\n not implemented", funcName);
    char *a = 0;
    *a = 0;
}

typedef struct {
    int64_t low;
    int64_t high;
} int128;

typedef struct {
    uint64_t low;
    uint64_t high;
} uint128;

typedef struct {
    double low;
    double high;
} LongDouble;

typedef uint64_t capi_name;
typedef char bool;

int32_t db_store_i64(uint64_t scope, capi_name table, capi_name payer, uint64_t id,  const char* data, uint32_t len)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void db_update_i64(int32_t iterator, capi_name payer, const char* data, uint32_t len)
{
    go_panic((char *)__FUNCTION__);
    return;
}

void db_remove_i64(int32_t iterator)
{
    go_panic((char *)__FUNCTION__);
    return;
}

int32_t db_get_i64(int32_t iterator, const char* data, uint32_t len)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_next_i64(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_previous_i64(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_find_i64(capi_name code, uint64_t scope, capi_name table, uint64_t id)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_lowerbound_i64(capi_name code, uint64_t scope, capi_name table, uint64_t id)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_upperbound_i64(capi_name code, uint64_t scope, capi_name table, uint64_t id)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_end_i64(capi_name code, uint64_t scope, capi_name table)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx64_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const uint64_t* secondary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void db_idx64_update(int32_t iterator, capi_name payer, const uint64_t* secondary)
{
    go_panic((char *)__FUNCTION__);
    return;
}

void db_idx64_remove(int32_t iterator)
{
    go_panic((char *)__FUNCTION__);
    return;
}

int32_t db_idx64_next(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx64_previous(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx64_find_primary(capi_name code, uint64_t scope, capi_name table, uint64_t* secondary, uint64_t primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx64_find_secondary(capi_name code, uint64_t scope, capi_name table, const uint64_t* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx64_lowerbound(capi_name code, uint64_t scope, capi_name table, uint64_t* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx64_upperbound(capi_name code, uint64_t scope, capi_name table, uint64_t* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx64_end(capi_name code, uint64_t scope, capi_name table)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx128_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const uint128* secondary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void db_idx128_update(int32_t iterator, capi_name payer, const uint128* secondary)
{
    go_panic((char *)__FUNCTION__);
    return;
}

void db_idx128_remove(int32_t iterator)
{
    go_panic((char *)__FUNCTION__);
    return;
}

int32_t db_idx128_next(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx128_previous(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx128_find_primary(capi_name code, uint64_t scope, capi_name table, uint128* secondary, uint64_t primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx128_find_secondary(capi_name code, uint64_t scope, capi_name table, const uint128* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx128_lowerbound(capi_name code, uint64_t scope, capi_name table, uint128* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx128_upperbound(capi_name code, uint64_t scope, capi_name table, uint128* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx128_end(capi_name code, uint64_t scope, capi_name table)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx256_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const uint128* data, uint32_t data_len )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void db_idx256_update(int32_t iterator, capi_name payer, const uint128* data, uint32_t data_len)
{
    go_panic((char *)__FUNCTION__);
    return;
}

void db_idx256_remove(int32_t iterator)
{
    go_panic((char *)__FUNCTION__);
    return;
}

int32_t db_idx256_next(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx256_previous(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx256_find_primary(capi_name code, uint64_t scope, capi_name table, uint128* data, uint32_t data_len, uint64_t primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx256_find_secondary(capi_name code, uint64_t scope, capi_name table, const uint128* data, uint32_t data_len, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx256_lowerbound(capi_name code, uint64_t scope, capi_name table, uint128* data, uint32_t data_len, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx256_upperbound(capi_name code, uint64_t scope, capi_name table, uint128* data, uint32_t data_len, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx256_end(capi_name code, uint64_t scope, capi_name table)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_double_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const double* secondary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void db_idx_double_update(int32_t iterator, capi_name payer, const double* secondary)
{
    go_panic((char *)__FUNCTION__);
    return;
}

void db_idx_double_remove(int32_t iterator)
{
    go_panic((char *)__FUNCTION__);
    return;
}

int32_t db_idx_double_next(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_double_previous(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_double_find_primary(capi_name code, uint64_t scope, capi_name table, double* secondary, uint64_t primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_double_find_secondary(capi_name code, uint64_t scope, capi_name table, const double* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_double_lowerbound(capi_name code, uint64_t scope, capi_name table, double* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_double_upperbound(capi_name code, uint64_t scope, capi_name table, double* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_double_end(capi_name code, uint64_t scope, capi_name table)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_long_double_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const LongDouble* secondary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void db_idx_long_double_update(int32_t iterator, capi_name payer, const LongDouble* secondary)
{
    go_panic((char *)__FUNCTION__);
    return;
}

void db_idx_long_double_remove(int32_t iterator)
{
    go_panic((char *)__FUNCTION__);
    return;
}

int32_t db_idx_long_double_next(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_long_double_previous(int32_t iterator, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_long_double_find_primary(capi_name code, uint64_t scope, capi_name table, LongDouble* secondary, uint64_t primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_long_double_find_secondary(capi_name code, uint64_t scope, capi_name table, const LongDouble* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_long_double_lowerbound(capi_name code, uint64_t scope, capi_name table, LongDouble* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_long_double_upperbound(capi_name code, uint64_t scope, capi_name table, LongDouble* secondary, uint64_t* primary)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t db_idx_long_double_end(capi_name code, uint64_t scope, capi_name table)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int64_t kv_erase(uint64_t contract, const char* key, uint32_t key_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int64_t kv_set(uint64_t contract, const char* key, uint32_t key_size, const char* value, uint32_t value_size, uint64_t payer)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

bool kv_get(uint64_t contract, const char* key, uint32_t key_size, uint32_t* value_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

uint32_t kv_get_data(uint32_t offset, char* data, uint32_t data_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

uint32_t kv_it_create(uint64_t contract, const char* prefix, uint32_t size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void kv_it_destroy(uint32_t itr)
{
    go_panic((char *)__FUNCTION__);
}

int32_t kv_it_status(uint32_t itr)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_compare(uint32_t itr_a, uint32_t itr_b)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_key_compare(uint32_t itr, const char* key, uint32_t size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_move_to_end(uint32_t itr)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_next(uint32_t itr, uint32_t* found_key_size, uint32_t* found_value_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_prev(uint32_t itr, uint32_t* found_key_size, uint32_t* found_value_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_lower_bound(uint32_t itr, const char* key, uint32_t size, uint32_t* found_key_size, uint32_t* found_value_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_key(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t* actual_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int32_t kv_it_value(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t* actual_size)
{
    go_panic((char *)__FUNCTION__);
    return 0;
}
