//go:build !tinygo.wasm
// +build !tinygo.wasm

#include "stacktrace.h"

extern "C" void print_stack_trace() {
    print_stacktrace();
}
