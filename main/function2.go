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
} Result;

// 定义函数指针类型
typedef Result (*create_person_func)(int, char*, char*);
typedef void (*free_person_func)(Result);

// 调用 create_person_func 的函数
Result call_create_person(create_person_func f, int id, char name[], char data[]) {
	return f(id, name, data);
}

// 调用 free_person_func 的函数
void call_free_person(free_person_func f, Result p) {
	f(p);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const config = "/home/lib/"

// 定义与 C 函数指针对应的 Go 类型
type CreatePersonFunc C.create_person_func
type FreePersonFunc C.free_person_func

var libHandles = make(map[string]func())
var handleMap = make(map[string]CreatePersonFunc)
var freeHandleMap = make(map[string]FreePersonFunc)

// 定义与C结构体对应的Go结构体
type Person struct {
	ID   int
	Name string
	Data []byte
}

func Handle2() {
	loadSo("face", "1", "", "")
	loadSo("face", "2", "", "")
	loadSo("face", "3", "", "")

	Handler(1, "face", "1", "", "")
	Handler(2, "face", "2", "", "")
	Handler(3, "face", "3", "", "")

	DestroyAllSo()
}

func callCreatePerson(createPerson CreatePersonFunc, id int, name string, data string) C.Result {
	cName := C.CString(name)
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cData))
	return C.call_create_person(createPerson, C.int(id), cName, cData)
}

// 包装 C 函数 call_free_person
func callFreePerson(freePerson FreePersonFunc, person C.Result) {
	C.call_free_person(freePerson, person)
}

func loadSo(name, version, base64, rmtp string) {
	sopath := config + name + "_" + version + ".so"
	key := name + "_" + version

	// Open the shared library
	libHandle := C.dlopen(C.CString(sopath), C.RTLD_LAZY)
	if libHandle == nil {
		fmt.Println("Error opening library")
		return
	}
	libHandles[key] = func() {
		// 回调函数释放库
		C.dlclose(libHandle)
		fmt.Println("Unloaded library success:", key)
	}

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

	handleMap[key] = createPerson
	freeHandleMap[key] = freePerson
	fmt.Println("Loaded library finish:", sopath)
}

func DestroySo(name, version string) {
	key := name + "_" + version
	handle, ok := libHandles[key]
	if !ok {
		fmt.Println("Library not found:", key)
		return
	}
	handle()
}

func DestroyAllSo() {
	for key, handle := range libHandles {
		// Close the shared library
		handle()

		// Remove the library handle from the map
		delete(libHandles, key)

		fmt.Println("Unloaded library:", key)
	}
}

func Handler(id int, name, version, base64, rmtp string) {
	key := name + "_" + version
	handleTask, ok := handleMap[key]
	if !ok {
		fmt.Println("Library not found:", key)
		return
	}

	freeHandle, ok := freeHandleMap[key]
	if !ok {
		fmt.Println("Library not found:", key)
		return
	}

	person := callCreatePerson(handleTask, id, name, version)

	p := Person{
		ID:   int(person.id),
		Name: C.GoString(person.name),
		Data: []byte(C.GoString(person.data)),
	}

	fmt.Println(p)
	fmt.Println(C.GoString(person.data))

	defer callFreePerson(freeHandle, person)

}
