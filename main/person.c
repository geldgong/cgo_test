#include <stdlib.h>
#include <string.h>

typedef struct {
    int id;
    char* name;
    char* data;
} Person;

Person create_person(int id, char name[], char data[]) {
    Person p;
    p.id = id;
    p.name = strdup(name); // 使用 strdup 分配内存并复制字符串
    p.data = strdup(data); // 使用 strdup 分配内存并复制字符串

    return p;
}

void free_person(Person p) {
    free(p.name);
    free(p.data);
}
