package main

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>
#include <string.h>

typedef struct {
    int id;
    char name[50];
    char data[100];
} Person;

Person create_person(int id, char name[], char data[]) {
    Person p;
    p.id = id;
    p.name = name;
    p.data = data;

    return p;
}
*/
import "C"
import "fmt"

func main() {
	cName := C.CString("Alice")
	p := C.create_person(1, cName, cName)
	fmt.Println(p)
}
