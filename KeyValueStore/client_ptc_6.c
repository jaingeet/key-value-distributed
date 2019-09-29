#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include "keyvalue.h"


// 100% write - Put number_of_keys using a single client
int main() {
//    printf("To kill a server and check how long it takes a server to revive:\n");

    int number_of_keys = 100;

    char *serverList[] = {
       "localhost:8001",
       "localhost:8003",
       NULL
    };

    time_t start_time, end_time;
    char* oldValue = malloc(1024);
    start_time = time(0);
    printf("Calling init %d \n", kv739_init(serverList, 2));
    kv739_put("123", "123", oldValue);
    int count = 0;

    for(int i = 0; i < number_of_keys; i++) {
        int start = i*5;
        int end = start + 5;

        // write 5 key-values
        for(int j = start; j<end; j++) {
            count++;
            printf("%d\n", j);
            char str[12];
            sprintf(str, "%d", j);
            kv739_put(str, str, oldValue);
        }

        // Read a key for 100 times
        for(int j = 0; j<95; j++) {
            count++;
            char str[12];
            sprintf(str, "%d", start);
            kv739_get(str, oldValue);
            printf("Get key for %d, %s\n", start, oldValue);
        }

    }
    printf("counter is : %d\n", count);
    end_time = time(0);

    double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    double throughput = (number_of_keys*100)/time_elapsed;
    double latency = time_elapsed/(number_of_keys*100);

    printf("throughput, latency is : %f, %f\n", (double) throughput, (double) latency);
}
