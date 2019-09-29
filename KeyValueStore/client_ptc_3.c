#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "keyvalue.h"

// 100% read - Get number_of_keys using a single client
int main() {
//    printf("To kill a server and check how long it takes a server to revive:\n");

    int number_of_keys = 1000;

    char *serverList[] = {
       "localhost:8001",
       "localhost:8002",
       "localhost:8003",
       NULL
    };

    time_t start_time, end_time;
    char* oldValue = malloc(1024);
    printf("Calling init %d \n", kv739_init(serverList, 3));

    // Put number_of_keys
    for(int i = 0; i < number_of_keys; i++) {
        char str[12];
        sprintf(str, "%d", i);
        kv739_put(str, str, oldValue);
    }

    // Get number_of_keys
    start_time = time(0);
    for(int i = 0; i < number_of_keys; i++) {
        char str[12];
        sprintf(str, "%d", i);
        kv739_get(str, oldValue);
    }
    end_time = time(0);

    double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    double throughput = number_of_keys/time_elapsed;
    double latency = time_elapsed/number_of_keys;

    printf("throughput, latency is : %f, %f\n", (double) throughput, (double) latency);
}
