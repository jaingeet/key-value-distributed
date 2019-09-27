#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "keyvalue.h"

// To kill a server and check how long it takes a server to revive
//TODO: update the PID for server 2 here in order to kill it

int main() {
    printf("To kill a server and check how long it takes a server to revive:\n");

    char* oldValue = malloc(1024);
    char* key = "523";
    char* kill_pid = "kill -9 4849"; //This should be the PID of server 3

    char *serverList[] = {
        "localhost:8003",
        NULL
    };

    time_t start_time, end_time;

    system(kill_pid);

    start_time = time(0);

    while(kv739_init(serverList, 1) == -1) {}

    end_time = time(0);

    double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
    printf("Calling shutdown  %d\n", kv739_shutdown());
}
