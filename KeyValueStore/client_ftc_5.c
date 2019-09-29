#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include <pthread.h>
#include "keyvalue.h"

// To test if all servers are in a sync state if one of the servers goes down 
// in the middle and then is revived automatically after sometime

//TODO: in script empty the txt files restart the server and pass server 1 index PID here

struct arg_struct {
    int start_key;
    int end_key;
    char *server_address;
};

char *serverList[] = {
    "localhost:8001",
    "localhost:8002",
    "localhost:8003",
    NULL
};

// 100% write - Put number_of_keys using a multiple client

// The function to be executed by all threads
void *putKeys(void *arguments)
{
    struct arg_struct *args = arguments;
    int start_key = args->start_key;
    int end_key = args->end_key;
    printf("start, end is %d, %d\n", start_key, end_key);

    char *server_list[] = {
        args->server_address,
        NULL
    };

    char* oldValue = malloc(1024);

    printf("Calling init %d \n", kv739_init(server_list, 1));

    for(int i = start_key; i < end_key; i++) {
        printf("key is, %d\n", i);
        char str[12];
        sprintf(str, "%d", i);
        kv739_put(str, str, oldValue);
        // if (i == 500) {
        //     char* kill_pid = "kill -9 3564";
        //     system(kill_pid);
        // }
    }
    return NULL;
}

int main()
{
	int num_of_threads = 3;
	int number_of_keys = 10000;
	int start_key = 0;
	int per_client_keys = number_of_keys/num_of_threads;
	time_t start_time, end_time;
    start_time = time(0);
    pthread_t threads[num_of_threads];
    struct arg_struct args[num_of_threads];

	// Let us create three threads
	for (int i = 0; i < num_of_threads; i++) {
	    args[i].start_key = start_key;
        args[i].end_key = start_key + per_client_keys;
        args[i].server_address = serverList[i];
	    printf("start, per_client_keys is %d, %d\n", args[i].start_key, args[i].end_key);
		pthread_create(&threads[i], NULL, putKeys, (void *)&args[i]);
		start_key += per_client_keys;
	}

	for(int i = 0; i<num_of_threads; i++) {
	    pthread_join(threads[i], NULL);
	}

	end_time = time(0);
	double time_elapsed = difftime(end_time, start_time);
    printf("start time, end time, timediff is : %f, %f, %f \n", (double) start_time, (double) end_time, time_elapsed);
//	pthread_exit(NULL);
	return 0;
}
