#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "./../keyvalue.h"

// Read 10% of the keys - Hot  keys
int main() {
    int number_of_keys = 10000;
    char* oldValue = malloc(1024);
    char *serverList[] = {
       "localhost:8001",
       "localhost:8002",
       "localhost:8003",
       NULL
    };

    printf("Calling init %d \n", kv739_init(serverList, 3));

    for(int i = 0; i < number_of_keys; i++) {
        char str[12];
        sprintf(str, "%d", i);
        kv739_put(str, str, oldValue);
    }

    time_t start_time, end_time;
    start_time = time(0);
    for(int i = 0; i < number_of_keys; i++) {
        int random_number = rand() % 1000 + 8000;
        char str[12];
        sprintf(str, "%d", random_number);
        kv739_get(str, oldValue);
    }
    end_time = time(0);

    double time_elapsed = difftime(end_time, start_time);
    double throughput = number_of_keys/time_elapsed;
    double latency = time_elapsed/number_of_keys;

    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    printf("throughput, latency is : %f, %f\n", (double) throughput, (double) latency);

    return 0;
}
