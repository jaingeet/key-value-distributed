package main

import (
	"log"
	"net/rpc"
)

type KeyValue struct {
	Key, Value string
}

var err error
var client Client
var reply string

func kv739_init(char** server_list) int {
	server_list="localhost:1234";
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	return 0;
}

func kv739_shutdown(void) int {
	return 0;
}

func kv739_get(char* key, char* value) int {
	client.Call("Task.GetKey", key, value);
	return 0;
}

func kv739_put(char* key, char* value, char* old_value) int {
	client.Call("Task.GetKey", key, value, old_value);
	return 0;
}

func main() {
}
