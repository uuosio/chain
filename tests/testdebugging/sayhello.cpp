//go:build !tinygo.wasm
// +build !tinygo.wasm

#include <stdio.h>
extern "C" void say_hello() {
    // printf("++++++hello,world\n");
}
