package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var config = []map[string]string{
	{
		"port":     "8001",
		"host":     "localhost",
		"filename": "./0.txt",
	},
	{
		"port":     "8002",
		"host":     "localhost",
		"filename": "./1.txt",
	},
	{
		"port":     "8003",
		"host":     "localhost",
		"filename": "./2.txt",
	},
}
var serverIndex int

var filename string

// Make a new KeyValue type that is a typed collection of fields
// (Key and Value), both of which are of type string
type KeyValue struct {
	Key, Value, TimeStamp string
}

type KeyValuePair struct {
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
	if scanner.Scan() {
		scanner.Text()
	}

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

//PutKey ...  TODO: can change timestamp type to Time instead of string
func (t *Task) PutKey(keyValue KeyValuePair, oldValue *string) error {
	//get file and find the given key
	//initialise the epoch timestamp
	// first line of the file should contain last updated timestamp
	// update last updated timestamp and the key timestamp
	// if EOF reached append new key, value and timestamp in the end of the file
	// send async requests to update the key value pair with timestamp in other replicas (if no response received
	// in callback from the other server, reinit the server that is not up)

	keyFound := false
	filePath := config[serverIndex]["filename"]
	curTimeStamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	newKeyValueString := string(keyValue.Key + "," + keyValue.Value + "," + curTimeStamp)

	fileContent, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
		return err
	}

	lines := strings.Split(string(fileContent), "\n")

	lines[0] = curTimeStamp

	for i := 1; i < len(lines); i++ {
		line := strings.Split(lines[i], ",")
		if line[0] == keyValue.Key {
			*oldValue = line[1]
			lines[i] = newKeyValueString
			keyFound = true
			break
		}
	}

	if !keyFound {
		lines = append(lines, newKeyValueString)
	}

	newFileContent := strings.Join(lines[:], "\n")
	err = ioutil.WriteFile(filename, []byte(newFileContent), 0)
	if err != nil {
		fmt.Printf("%s ", err)
		return err
	}

	for index := range config {
		if index != serverIndex {
			client, err := rpc.DialHTTP("tcp", config[index]["host"]+":"+config[index]["port"])
			if err != nil {
				// callback (check reply and restart server if there is a connection error)
				fmt.Printf("%s ", err)
				RestartServer(index)
			} else {
				client.Go("Task.SyncKey", KeyValue{Key: keyValue.Key, Value: keyValue.Value, TimeStamp: curTimeStamp}, nil, nil)
			}
		}
	}

	return nil
}

func RestartServer(serverIndex int) {
	// here -r is for server restart
	cmd := exec.Command("go", "run", "server.go", strconv.Itoa(serverIndex), " -r", " &")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("error\n")
		log.Fatal(err)
	}
	pid := cmd.Process.Pid
	fmt.Printf("Server %d restarts with process id: %d\n", serverIndex, pid)
}

func SyncReplicas(time int64) error {
	// get all the key value and timestamp pairs from other servers that were updated after this timestamp
	// call synckey method for all the key value pairs

	var updates []KeyValue

	//TODO: make a connection and an RPC Call here? (also need to get updates from all servers)
	for index, _ := range config {
		if index != serverIndex {
			client, err := rpc.DialHTTP("tcp", config[index]["host"]+":"+config[index]["port"])
			if err != nil {
				log.Fatal(err)
			} else {
				err := client.Go("Task.GetUpdates", time, &updates, nil)
				if err != nil {
					RestartServer(index)
				}
			}
		}
	}

	// To Do - add functionality for bulk update
	for _, update := range updates {
		SyncKeyLocally(update)
	}

	return nil
}

//GetUpdates ... to get key values updated after the timestamp
func (t *Task) GetUpdates(timestamp string, updates *[]KeyValue) error {
	// return an array of all key,value and timestamp where timestamp > given timestamp

	file, err := os.Open(config[serverIndex]["filename"])

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Ignore the first line as it contains the last updated time

	if scanner.Scan() {
		scanner.Text()
	}

	for scanner.Scan() {
		var line []string = strings.Split(scanner.Text(), ",")
		if line[2] >= timestamp {
			*updates = append(*updates, KeyValue{Key: line[0], Value: line[1], TimeStamp: line[2]})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func SyncKeyLocally(keyValue KeyValue) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}

	newKeyValueString := keyValue.Key + "," + keyValue.Value + "," + keyValue.TimeStamp
	lines := strings.Split(string(data), "\n")
	var found = false // to check if the key is present or not

	var TimeInFile int64 = 0
	UpdatedTime, _ := strconv.ParseInt(keyValue.TimeStamp, 10, 64)
	for i, line := range lines {
		var array []string = strings.Split(line, ",")
		if i != 0 {
			if keyValue.Key == array[0] {
				TimeInFile, err = strconv.ParseInt(array[2], 10, 64)
				if TimeInFile < UpdatedTime {
					lines[i] = newKeyValueString
				}
				found = true
				break
			}
		}
	}

	TimeInFile, err = strconv.ParseInt(lines[0], 10, 64)

	if TimeInFile < UpdatedTime {
		lines[0] = strconv.FormatInt(UpdatedTime, 10)
	}

	fmt.Println()

	if !found {
		lines = append(lines, newKeyValueString)
	}

	newFileContent := strings.Join(lines[:], "\n")
	err = ioutil.WriteFile(filename, []byte(newFileContent), 0)

	if err != nil {
		fmt.Printf("%s ", err)
		return err
	}

	return nil
}

func (t *Task) SyncKey(keyValue KeyValue, reply *string) error {
	// find key in the file
	// if found: check the updated timestamp for the key
	// if updated timestamp > timestamp arg (do nothing)
	// if equal ? (CASE NEEDS TO BE GIVEN A THOUGHT)
	// else: update the key with the passed value and timestamp
	return SyncKeyLocally(keyValue)
}

//Init ... takes in config and index of the current server in config
func Init(index int, restart bool) error {
	// create file and add first line if not already present
	// sync after all servers are up
	task := new(Task)
	// Publish the receivers methods
	err := rpc.Register(task)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}

	// Sync with the other server before restart
	var time int64 = 0
	if restart == true {
		// read the last updated timestamp (-5 seconds etc -- for failover delay) from the file
		file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		var timestamp string = scanner.Text()
		time, err = strconv.ParseInt(timestamp, 10, 64)
		time = time - 5*1000
	}

	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", config[index]["host"]+":"+config[index]["port"])
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", config[index]["port"])
	// Start accept incoming HTTP connections
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
	if restart == true {
		SyncReplicas(time)
	}
	return nil
}

func main() {
	args := os.Args[1:]
	fmt.Printf("len %d", len(args))
	serverIndex, _ = strconv.Atoi(args[0])
	filename = config[serverIndex]["filename"]
	var restart bool = false
	if len(args) > 2 {
		restart = true
	}
	Init(serverIndex, restart)
}
