package main

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>
#include <string.h>

typedef struct {
    int id;
    char* name;
    char* data;
} Person;

// 定义函数指针类型
typedef Person (*create_person_func)(int, char*, char*);
typedef void (*free_person_func)(Person);

// 调用 create_person_func 的函数
Person call_create_person(create_person_func f, int id, char name[], char data[]) {
	return f(id, name, data);
}

// 调用 free_person_func 的函数
void call_free_person(free_person_func f, Person p) {
	f(p);
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

// 定义与C结构体对应的Go结构体
type Person1 struct {
	ID   int
	Name string
	Data []byte
}

// 包装 C 函数 call_create_person
func callCreatePerson1(createPerson CreatePersonFunc, id int, name string, data string) C.Person {
	cName := C.CString(name)
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cData))
	return C.call_create_person(createPerson, C.int(id), cName, cData)
}

// 包装 C 函数 call_free_person
func callFreePerson1(freePerson FreePersonFunc, person C.Person) {
	C.call_free_person(freePerson, person)
}

func Handle3() {
	//cName := C.CString("Alice")
	//cData := C.CString("Some data")
	//defer C.free(unsafe.Pointer(cName))
	//defer C.free(unsafe.Pointer(cData))
	//
	//p := C.create_person(1, cName, cData)
	//defer C.free_person(p) // 释放 C 分配的内存
	//
	//// 将 C 字符串转换为 Go 字符串
	//name := C.GoString(p.name)
	//data := C.GoString(p.data)
	//
	//fmt.Println("ID:", p.id)
	//fmt.Println("Name:", name)
	//fmt.Println("Data:", data)

	// Path to the shared library
	soPath := "./libperson.so"

	// Open the shared library
	libHandle := C.dlopen(C.CString(soPath), C.RTLD_LAZY)
	if libHandle == nil {
		fmt.Println("Error opening library")
		return
	}
	defer C.dlclose(libHandle)

	// Load the function symbols
	createPersonSymbol := C.dlsym(libHandle, C.CString("create_person"))
	if createPersonSymbol == nil {
		fmt.Println("Error loading symbol 'create_person'")
		return
	}

	freePersonSymbol := C.dlsym(libHandle, C.CString("free_person"))
	if freePersonSymbol == nil {
		fmt.Println("Error loading symbol 'free_person'")
		return
	}

	// Cast the symbols to the appropriate function types
	createPersonFunc := (*[0]byte)(unsafe.Pointer(createPersonSymbol))
	freePersonFunc := (*[0]byte)(unsafe.Pointer(freePersonSymbol))

	// Convert the function pointers
	createPerson := *(*CreatePersonFunc)(unsafe.Pointer(&createPersonFunc))
	freePerson := *(*FreePersonFunc)(unsafe.Pointer(&freePersonFunc))

	// Prepare the data to pass to the C function
	cName := C.CString("Alice")
	cData := C.CString("Some data")
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cData))

	person := callCreatePerson1(createPerson, 1, "Alice", "Some data")
	defer callFreePerson1(freePerson, person)
	fmt.Println("ID:", person.id)
	fmt.Println("Name:", C.GoString(person.name))
	fmt.Println("Data:", C.GoString(person.data))

	p := Person1{
		ID:   int(person.id),
		Name: C.GoString(person.name),
		Data: []byte(C.GoString(person.data)),
	}

	fmt.Println(p)
}
