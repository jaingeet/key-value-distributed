#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include "./../keyvalue.h"

// 100% read - Get number_of_keys using a single client
int main() {
    int number_of_keys = 10000;
    char* oldValue = malloc(1024);
    char *serverList[] = {
       "10.10.1.1:8001",
       "10.10.1.2:8002",
       "10.10.1.3:8003",
       NULL
    };

    printf("Calling init %d \n", kv739_init(serverList, 3));

    // Put number_of_keys
    for(int i = 0; i < number_of_keys; i++) {
        char str[12];
        // sprintf(str, "%d", i);
        kv739_put(str, str, oldValue);
        // printf("%d", i);
    }
    time_t start_time, end_time;

    int not_found = 0;
    int wrong_values_count = 0;

    // Get number_of_keys
    start_time = time(0);
    for(int i = 0; i < number_of_keys; i++) {
        char str[12];
        // sprintf(str, "%d", i);
        int x = kv739_get(str, oldValue);
        if( x == -1) {
            not_found++;
        }
        if(strcmp(str, oldValue) != 0) {
            // printf("wrong => %s, %s\n", str, oldValue);
            wrong_values_count++;
        }
    }
    end_time = time(0);

    printf("Keys Not Found => %d\n", not_found);
    printf("Value wrong Found => %d\n", wrong_values_count);

    double time_elapsed = difftime(end_time, start_time);
    double throughput = number_of_keys/time_elapsed;
    double latency = time_elapsed/number_of_keys;

    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    printf("throughput, latency is : %f, %f\n", (double) throughput, (double) latency);

    return 0;
}
