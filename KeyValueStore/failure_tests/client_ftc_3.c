#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "keyvalue.h"

// To test the time taken for the key update to propagate to other server if the server on which the operation happened fails immediately after
//TODO: delete contents of files after first line, update the PID for server 2 here in order to kill it


// clien1 put x server2
// kill server process (read from file -- line 1)
// get x from server 0 and server 1
// get from server 2 (if server 2 was up now)
// pass the pid of server 2 in this file

int main() {
    printf("To test the time taken for the key update to propagate to other server if the server on which the operation happened fails immediately after:\n");

    char* oldValue = malloc(1024);
    char* key = "1023";
    char* kill_pid = "kill -9 3564"; //This should be the PID of server 3

    char *serverList[] = {
        "localhost:8003",
        NULL
    };

    printf("Calling init %d \n", kv739_init(serverList, 1));

    printf("Calling Put %d \n", kv739_put("a", key, oldValue));

    system(kill_pid);

    printf("Returned old value is %s \n", oldValue);
    printf("Calling shutdown %d\n", kv739_shutdown());

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

// How to bring the server up if it fails? -- It should come up automatically
// Have written a script to do that