#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "keyvalue.h"
#include<string.h>

// To test the time taken for the key update to propagate to another server
//TODO: in script empty the files from second line onwards


// client1: connect to server 0
// make put requests
// shutdown connection

// after init with another server
// start time
// for loop (make get request from server 1) unless the key is same as the put value
// when it matches becomes the end time
// and then measure the time diff

int main() {
    printf("To test the time taken for the key update to propagate to another server:\n");

    char* oldValue = malloc(1024);
    char* key = "123";

    char *serverList[] = {
        "localhost:8001",
        NULL
    };

    printf("Calling init %d \n", kv739_init(serverList, 1));

    printf("Calling Put %d \n", kv739_put("a", key, oldValue));
    printf("Returned old value is %s \n", oldValue);

    printf("Calling shutdown  %d\n", kv739_shutdown());

    *serverList = "localhost:8002";

    printf("Calling init %d \n", kv739_init(serverList, 1));

    *oldValue = 0;
    
    time_t start_time, end_time;

    start_time = time(0);

    while(strcmp(oldValue, key) != 0) {
        kv739_get("a", oldValue);
    }
    end_time = time(0);
    double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    printf("Calling shutdown  %d\n", kv739_shutdown());
}
