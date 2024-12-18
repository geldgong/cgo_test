package main

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>

typedef void (*hello_func)();
typedef int (*add_func)(int, int);

void call_hello(hello_func f) {
    f();
}

int call_add(add_func f, int a, int b) {
    return f(a, b);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func Handle1() {
	// Path to the shared library
	soPath := "./libexample.so"

	// Open the shared library
	libHandle := C.dlopen(C.CString(soPath), C.RTLD_LAZY)
	if libHandle == nil {
		fmt.Println("Error opening library")
		return
	}
	defer C.dlclose(libHandle)

	// Load the function symbols
	helloSymbol := C.dlsym(libHandle, C.CString("hello_from_c"))
	if helloSymbol == nil {
		fmt.Println("Error loading symbol 'hello_from_c'")
		return
	}

	addSymbol := C.dlsym(libHandle, C.CString("add"))
	if addSymbol == nil {
		fmt.Println("Error loading symbol 'add'")
		return
	}

	// Cast the symbols to the appropriate function types
	helloFunc := (*C.hello_func)(unsafe.Pointer(helloSymbol))
	addFunc := (*C.add_func)(unsafe.Pointer(addSymbol))

	// Call the functions
	C.call_hello(C.hello_func(unsafe.Pointer(helloFunc)))

	result := C.call_add(C.add_func(unsafe.Pointer(addFunc)), 3, 4)
	fmt.Printf("Result of add(3, 4): %d\n", result)
}
