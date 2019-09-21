package main

import "C"

import (
	"log"
	"math/rand"
	"net/rpc"
)

var err error
var client Client
var reply string

//export kv739_init
func kv739_init(serverList []string) int {
	rand := rand.Intn(len(serverList))
	address := serverList[rand]
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("Connection error: ", err)
		return -1
	}
	return 0
}

//export kv739_shutdown
func kv739_shutdown() int {
	err := client.Close()
	if err != nil {
		log.Fatal("Unable to shutdown client connection: ", err)
		return -1
	}
	return 0
}

//export kv739_get
func kv739_get(key string, value string) int {
	err := client.Call("Task.GetKey", key, &value)
	if err != nil {
		log.Fatal("Could not get key: ", key, err)
		return -1
	}
	if len(value) > 1 {
		return 0
	}
	return 1
}

//export kv739_put
func kv739_put(key, value, oldValue string) int {
	err := client.Call("Task.GetKey", key, value, &oldValue)
	if err != nil {
		log.Fatal("Could not put key: ", key, " value: ", value, err)
		return -1
	}

	if len(oldValue) > 1 {
		return 0
	}
	return 1
}

func main() {
}
