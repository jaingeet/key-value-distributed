#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "./../keyvalue.h"

char *serverList[] = {
    "10.10.1.3:8003",
    NULL
};

// 33% write - using third client
int main()
{
    char* oldValue = malloc(1024);
	int start_key = 6666;
    int end_key = 9999;

    printf("Calling init %d \n", kv739_init(serverList, 1));

    time_t start_time, end_time;
    start_time = time(0);
	for (int i = start_key; i < end_key; i++) {
        char str[12];
        // sprintf(str, "%d", i);
        kv739_put(str, str, oldValue);
	}
	end_time = time(0);

	double time_elapsed = difftime(end_time, start_time);
    double throughput = (end_key - start_key)/time_elapsed;
    double latency = time_elapsed/(end_key - start_key);

    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    printf("throughput, latency is : %f, %f\n", (double) throughput, (double) latency);

	return 0;
}
