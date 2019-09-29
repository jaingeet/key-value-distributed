#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include "./../keyvalue.h"

// 95% read and 5% write
int main() {
    int number_of_keys = 100;
    char* oldValue = malloc(1024);
    char *serverList[] = {
       "localhost:8001",
       "localhost:8002",
       "localhost:8003",
       NULL
    };

    printf("Calling init %d \n", kv739_init(serverList, 3));

    time_t start_time, end_time;
    start_time = time(0);

    for(int i = 0; i < number_of_keys; i++) {
        int start = i*5;
        int end = start + 5;

        // write 5 key-values
        for(int j = start; j<end; j++) {
            printf("%d\n", j);
            char str[12];
            sprintf(str, "%d", j);
            kv739_put(str, str, oldValue);
        }

        // Read a key for 95 times
        for(int j = 0; j<95; j++) {
            char str[12];
            sprintf(str, "%d", start);
            kv739_get(str, oldValue);
        }
    }
    end_time = time(0);

    double time_elapsed = difftime(end_time, start_time);
    double throughput = (number_of_keys*100)/time_elapsed;
    double latency = time_elapsed/(number_of_keys*100);

    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    printf("throughput, latency is : %f, %f\n", (double) throughput, (double) latency);

    return 0;
}
