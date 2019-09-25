#include <stdio.h>
#include <stdlib.h>
#include "keyvalue.h"

int main() {
    printf("Using keyvalue lib from C:\n");

    char *serverList[] = {
        "localhost:8001",
        "localhost:8002",
        "localhost:8003",
        NULL
    };

    printf("keyvalue.keyvalue_init(12,99) = %d\n", kv739_init(serverList, 3));

    char* oldValue = malloc(1024);
    printf("calling get function ---- %d \n", kv739_put("a", "523", oldValue));
    printf("old value for key a is === %s ", oldValue);
    printf("getting key value a from the server --- %d ", kv739_get("a", oldValue));
    printf("value for key a from server is === %s ", oldValue);
    printf("calling get function ---- %d \n", kv739_put("b", "524", oldValue));
    printf("old value for key b is === %s ", oldValue);
    printf("getting key value b from the server --- %d ", kv739_get("b", oldValue));
    printf("value for key b from server is === %s ", oldValue);
    printf("calling get function ---- %d \n", kv739_put("a", "523", oldValue));
    printf("old value for key b is === %s ", oldValue);
    printf("getting key value b from the server --- %d ", kv739_get("a", oldValue));
    printf("value for key b from server is === %s ", oldValue);


    printf("calling server shutdown = %d\n", kv739_shutdown());

    printf("this should return an error: getting key value a from the server --- %d ", kv739_get("a", oldValue));

    // getCharPtr(serverList);

    //Call Cosine() - passing float param, float returned
    // printf("awesome.Cosine(1) = %f\n", (float)(Cosine(1.0)));
    
    // //Call Sort() - passing an array pointer
    // GoInt data[6] = {77, 12, 5, 99, 28, 23};
    // GoSlice nums = {data, 6, 6};
    // Sort(nums);
    // printf("awesome.Sort(77,12,5,99,28,23): ");
    // for (int i = 0; i < 6; i++){
    //     printf("%d,", ((GoInt *)nums.data)[i]);
    // }
    // printf("\n");

    // //Call Log() - passing string value
    GoString msg = {"Hello from C!", 13};
    // Log(msg);
}