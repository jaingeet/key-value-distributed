package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var config []map[string]string
var serverIndex int

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
	// use cache (may be later)
	// find the key in file and return
	// return null if not found
	return nil
}

func (t *Task) PutKey(key string, value string, old_value *string) error {
	*old_value = "1234"

	//get file and find the given key

	//initialise the epoch timestamp

	// first line of the file should contain last updated timestamp

	// update last updated timestamp and the key timestamp

	// if EOF reached append new key, value and timestamp in the end of the file

	// send async requests to update the key value pair with timestamp in other replicas (if no response received
	// in callback from the other server, reinit the server that is not up)

	return nil
}

func (t *Task) SyncReplicas(reply *string) error {
	// read the last updated timestamp (-5 seconds etc -- for failover delay) from the file
	// get all the key value and timestamp pairs from other servers that were updated after this timestamp
	// call synckey method for all the key value pairs
	return nil
}

func (t *Task) GetUpdates(timestamp string, reply **string) error {
	// return an array of all key,value and timestamp where timestamp > given timestamp
	return nil
}

func (t *Task) SyncKey(key string, value string, timestamp string, reply *string) error {
	// find key in the file
	// if found: check the updated timestamp for the key
	// if updated timestamp > timestamp arg (do nothing)
	// if equal ? (CASE NEEDS TO BE GIVEN A THOUGHT)
	// else: update the key with the passed value and timestamp
	return nil
}

//Init ... takes in config and index of the current server in config
func Init(config []map[string]string, index int) {
	task := new(Task)
	// Publish the receivers methods
	err := rpc.Register(task)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", host+":"+port)
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
