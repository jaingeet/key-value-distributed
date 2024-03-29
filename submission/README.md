# key-value-distributed

Start 3 servers: They are currently running on ports 8001, 8002, 8003. Make sure the ports are not already in use.

```
./server 0
./server 1
./server 2
```

After starting the servers, please start the below process that will check if the servers are up and restart them if they are not. Make sure you start the servers before running this script otherwise this will automatically start your servers:

```
python3 process_checker_daemon.py
```

Please export LD_LIBRARY_PATH to the directory where keyvalue.so is present.

```
export LD_LIBRARY_PATH=~/submission/
```

To build client1.c use the below command:

```
gcc -o client ./client1.c ./keyvalue.so
```

Please Note: When creating a new client, our library expect an integer argument specifying the number of items in the server_list.

```
int kv739_init(char ** server_list, int num_servers)
```

Also we require that you create Server list to pass to init as follows. Please use only these 3 address. You may use it in any combination or length varying from 1-3.

```
char *serverList[] = {
    "localhost:8001",
    "localhost:8002",
    "localhost:8003",
    NULL
};
```

The value pointers returned from get and put functions should be allocated memory:

```
char* oldValue = malloc(1024);
```

To run the client:

```
./client
```

YAYYY.... Please contact us if you face any issues.

Email: geetika@wisc.edu, sanand@cs.wisc.edu

