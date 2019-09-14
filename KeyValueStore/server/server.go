package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Make a new KeyValue type that is a typed collection of fields
// (Key and Value), both of which are of type string
type KeyValue struct {
	Key, Value string
}

// TODO: Do we need this?
type EditKeyValue struct {
	Key, Value, OldValue string
}

type Task int

// This should be fetched from file (We need a persistent store)
var keyValueStore []KeyValue

// GetToDo takes a string type and returns a ToDo
func (t *Task) GetKey(key string, value *string) error {
	*value = "123"
	return nil
}

func (t *Task) PutKey(key string, value string, old_value *string) error {
	*old_value = "1234"
	return nil
}

func Init(host string, port string) {
	task := new(Task)
	// Publish the receivers methods
	err := rpc.Register(task)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", host + ":" + port)
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", 1234)
	// Start accept incoming HTTP connections
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}

func main() {
	
}
