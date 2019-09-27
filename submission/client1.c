#include <stdio.h>
#include <stdlib.h>
#include "keyvalue.h"

int main() {
    printf("Using keyvalue lib from C:\n");

    char* oldValue = malloc(1024);

    char *serverList[] = {
        "localhost:8001",
        "localhost:8002",
        "localhost:8003",
        NULL
    };

    printf("Calling init %d \n", kv739_init(serverList, 3));

    printf("Calling Put %d \n", kv739_put("a", "523", oldValue));
    printf("Returned old value is %s \n", oldValue);

    printf("Calling Get %d \n", kv739_get("a", oldValue));
    printf("Returned value is %s \n", oldValue);

    printf("Calling shutdown  %d\n", kv739_shutdown());
}
