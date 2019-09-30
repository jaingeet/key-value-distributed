#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "../keyvalue.h"


char *serverList[] = {
    "localhost:8001",
    // "localhost:8002",
    NULL
};

// To test if all servers are in a sync state if one of the servers goes down 
// in the middle and then is revived automatically after sometime

//TODO: in script empty the txt files restart the server and pass server 0 index PID here

int main()
{
    char* oldValue = malloc(1024);
	int start_key = 0;
    int end_key = 3333;
	time_t start_time, end_time;
    
    printf("Calling init %d \n", kv739_init(serverList, 1));
    
    start_time = time(0);
	for (int i = start_key; i < end_key; i++) {
        printf("key is, %d\n", i);
        char str[12];
        sprintf(str, "%d", i);
        kv739_put(str, str, oldValue);
        if (i == 500) {
            char* kill_pid = "kill -9 39421";
            system(kill_pid);
        }
	}
	end_time = time(0);

	double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
	return 0;
}
