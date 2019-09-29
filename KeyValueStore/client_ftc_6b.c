#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "keyvalue.h"


char *serverList[] = {
    "localhost:8002",
    // "localhost:8003",
    NULL
};

// 100% write - Put number_of_keys using a multiple client

int main()
{
    char* oldValue = malloc(1024);
	int start_key = 3333;
    int end_key = 6666;
	time_t start_time, end_time;
    
    printf("Calling init %d \n", kv739_init(serverList, 1));
    
    start_time = time(0);
	for (int i = start_key; i < end_key; i++) {
        printf("key is, %d\n", i);
        char str[12];
        sprintf(str, "%d", i);
        kv739_put(str, str, oldValue);
	}
	end_time = time(0);

	double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
	return 0;
}
