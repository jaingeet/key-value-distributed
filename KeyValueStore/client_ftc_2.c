#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "keyvalue.h"

// To test the time taken to propagate an update to the failed server when update happened
//TODO: delete contents of files after first line, shutdown server 0


// server 0 is down at start
// server1 puts x
// shutdown
// read x from server 0 (loop until you get the value)

int main() {
    printf("To test the time taken for the key update to propagate to another server:\n");

    char* oldValue = malloc(1024);
    char* key = "523";

    char *serverList[] = {
        "localhost:8002",
        NULL
    };

    printf("Calling init %d \n", kv739_init(serverList, 1));

    printf("Calling Put %d \n", kv739_put("a", key, oldValue));
    printf("Returned old value is %s \n", oldValue);

    printf("Calling shutdown %d\n", kv739_shutdown());

    // Sync replica restarts the server before this is called (so it's awesome)
    system("./server/server 0");

    *serverList = "localhost:8001";

    printf("Calling init %d \n", kv739_init(serverList, 1));

    *oldValue = 0;

    time_t start_time, end_time;

    start_time = time(0);
    while(*oldValue != *key) {
        kv739_get("a", oldValue);
    }
    end_time = time(0);
    double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    printf("Calling shutdown  %d\n", kv739_shutdown());
}

// After this test case:
// Server 2 is crashing: Returning index out of range error in function GetUpdates
// This was because an empty line was inserted in 2.txt
// Why was an empty line inserted?

