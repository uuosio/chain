#include <stddef.h>
#include <stdint.h>
#include <stdio.h>

extern void print_stack_trace();

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

typedef struct {
	uint64_t data[4];
} capi_checksum256;

typedef struct {
	uint64_t data[8];
} capi_checksum512;

typedef struct {
	unsigned char data[20];
} capi_checksum160;

static void go_panic(char* funcName) {
    print_stack_trace();
    printf("%s\n not implemented", funcName);
    char *a = 0;
    *a = 0;
}

uint32_t get_active_producers( capi_name* producers, uint32_t datalen )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int64_t get_permission_last_used( capi_name account, capi_name permission )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int64_t get_account_creation_time( capi_name account )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void get_resource_limits( capi_name account, int64_t* ram_bytes, int64_t* net_weight, int64_t* cpu_weight )
{
    go_panic((char *)__FUNCTION__);
    return;
}

void set_resource_limits( capi_name account, int64_t ram_bytes, int64_t net_weight, int64_t cpu_weight )
{
    go_panic((char *)__FUNCTION__);
    return;
}

int64_t set_proposed_producers( char *producer_data, uint32_t producer_data_size )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

int64_t set_proposed_producers_ex( uint64_t producer_data_format, char *producer_data, uint32_t producer_data_size )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

bool is_privileged( capi_name account )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void set_privileged( capi_name account, bool is_priv )
{
    go_panic((char *)__FUNCTION__);
    return;
}

void set_blockchain_parameters_packed( char* data, uint32_t datalen )
{
    go_panic((char *)__FUNCTION__);
    return;
}

uint32_t get_blockchain_parameters_packed( char* data, uint32_t datalen )
{
    go_panic((char *)__FUNCTION__);
    return 0;
}

void preactivate_feature( const capi_checksum256* feature_digest )
{
    go_panic((char *)__FUNCTION__);
    return;
}

int32_t
check_transaction_authorization( const char* trx_data,     uint32_t trx_size,
								const char* pubkeys_data, uint32_t pubkeys_size,
								const char* perms_data,   uint32_t perms_size
							);
int32_t check_permission_authorization( capi_name account,
								capi_name permission,
								const char* pubkeys_data, uint32_t pubkeys_size,
								const char* perms_data,   uint32_t perms_size,
								uint64_t delay_us
							) {
	go_panic((char *)__FUNCTION__);
	return 0;
}

int32_t check_transaction_authorization( const char* trx_data,     uint32_t trx_size,
								const char* pubkeys_data, uint32_t pubkeys_size,
								const char* perms_data,   uint32_t perms_size
							) {
	go_panic((char *)__FUNCTION__);
	return 0;
}
void set_action_return_value(char *return_value, size_t size) {
	go_panic((char *)__FUNCTION__);
}

void set_kv_parameters_packed( const char* data, uint32_t datalen ) {

}
