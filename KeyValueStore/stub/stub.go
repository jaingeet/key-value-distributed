package main

import (
	"log"
	"net/rpc"
)

const config []map[string]string

config[0] = map[string]string{
    "port": "0000",
    "host": "127.0.0.1",
    "filename": "0.txt",
}
config[1] = map[string]string{
    "port": "0001",
    "host": "127.0.0.1",
    "filename": "1.txt",
}
config[2] = map[string]string{
    "port": "0002",
    "host": "127.0.0.1",
    "filename": "2.txt",
}

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
