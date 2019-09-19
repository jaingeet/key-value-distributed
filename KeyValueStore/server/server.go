package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

var config []map[string]string
var serverIndex int


var filename string = config[serverIndex]["filename"]

// Make a new KeyValue type that is a typed collection of fields
// (Key and Value), both of which are of type string
type KeyValue struct {
	Key, Value, TimeStamp string
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
	// use cache (may be later)
	// find the key in file and return
	// return null if not found

	file, err := os.Open(config[serverIndex]["filename"])
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Ignore the first line as it contains the last updated time
	// Not sure if we also need to call Text() to move the pointer
	scanner.Scan()

	for scanner.Scan() {
		var line []string = strings.Split(scanner.Text(), ",")
		if line[0] == key {
			*value = line[1]
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}

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

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan();

	var timestamp string = scanner.Text()
	var time int64 = 0;
	time, err = strconv.ParseInt(timestamp, 10, 64)
	time = time - 5*1000;

	var updates = GetUpdates(time)

	// To Do - add functionality for bulk update
	for _, update := range updates {
		SyncKey(update.Key, update.Value, update.TimeStamp)
	}

	return nil
}

// We don't want to make this an RPC Call
func GetUpdates(timestamp string) []KeyValue {
	// return an array of all key,value and timestamp where timestamp > given timestamp

	var updates []KeyValue

	file, err := os.Open(config[serverIndex]["filename"])

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Ignore the first line as it contains the last updated time
	scanner.Scan()

	for scanner.Scan() {
		var line []string = strings.Split(scanner.Text(), ",")
		if line[2] >= timestamp {
			updates = append(updates, KeyValue{Key: line[0], Value: line[1], TimeStamp: line[2]})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return updates
}

func (t *Task) SyncKey(key string, value string, timestamp string, reply *string) error {
	// find key in the file
	// if found: check the updated timestamp for the key
	// if updated timestamp > timestamp arg (do nothing)
	// if equal ? (CASE NEEDS TO BE GIVEN A THOUGHT)
	// else: update the key with the passed value and timestamp

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}

	lines := strings.Split(string(data), "\n")
	var found = false // to check if the key is present or not

	var TimeInFile int64 = 0;
	var UpdatedTime int64 = 0;
	for i, line := range lines {
		var array []string = strings.Split(line, ",")
		if i != 0 {
			TimeInFile, err = strconv.ParseInt(array[2], 10, 64)
			UpdatedTime, err = strconv.ParseInt(timestamp, 10, 64)
			if key == array[0] {
				if TimeInFile < UpdatedTime {
					lines[i] = key + "," + value + "," + timestamp;
				}
				found = true;
				break;
			}
		}
	}

	TimeInFile, err = strconv.ParseInt(lines[0], 10, 64)
	UpdatedTime, err = strconv.ParseInt(timestamp, 10, 64)

	if TimeInFile < UpdatedTime {
		lines[0] = strconv.FormatInt(UpdatedTime, 10);
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// if key is not found, we need to append it to data
	if !found {
		file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if _, err := file.Write([]byte("\n" + key + "," + value + "," + timestamp)); err != nil {
			log.Fatal(err)
			return err
		}
		if err := file.Close(); err != nil {
			log.Fatal(err)
			return err
		}
	}

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
