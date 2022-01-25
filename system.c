#include <stdio.h>
#include <stdint.h>

void prints_l( const char* cstr, uint32_t len);

int read(int fd, void *buf, size_t count) {
    return 0;
}

int write(int fd, const void *buf, size_t count) {
    prints_l((const char*)buf, count);
    return count;
}

uint64_t lseek(int fd, off_t offset, int whence) {
    return 0;
}

int open(const char *pathname, int flags, int mode) {
    return -1;
}

int close(int fd) {
    return -1;
}
