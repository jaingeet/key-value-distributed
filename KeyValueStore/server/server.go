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
	"strconv"
	"strings"
	"sync"
	"time"
)

type DLinkedList struct {
	val  string
	key  string
	prev *DLinkedList
	next *DLinkedList
}

type LRUCache struct {
	capacity int
	head     *DLinkedList
	tail     *DLinkedList
	cache    map[string]*DLinkedList
}

func (lru *LRUCache) Init(cap int) {
	fmt.Println("calling constructor")
	lru.capacity = cap
	lru.cache = make(map[string]*DLinkedList)
}

func (c *LRUCache) Get(key string) string {
	l, ok := c.cache[key]
	if !ok {
		return ""
	}
	c.removeFromChain(l)
	c.addToChain(l)
	return l.val
}

func (c *LRUCache) Put(key string, value string) {
	l, ok := c.cache[key]
	if !ok {
		l = &DLinkedList{val: value, key: key}
		c.cache[key] = l
	} else {
		l.val = value
		c.removeFromChain(l)
	}

	c.addToChain(l)
	if len(c.cache) > c.capacity {
		l = c.tail
		if l != nil {
			c.removeFromChain(l)
			delete(c.cache, l.key)
		}
	}
}

func (c *LRUCache) addToChain(l *DLinkedList) {
	l.prev = nil
	l.next = c.head
	if l.next != nil {
		l.next.prev = l
	}
	c.head = l
	if c.tail == nil {
		c.tail = l
	}
}

func (c *LRUCache) removeFromChain(l *DLinkedList) {
	if l == c.head {
		c.head = l.next
	}

	if l == c.tail {
		c.tail = l.prev
	}

	if l.next != nil {
		l.next.prev = l.prev
	}

	if l.prev != nil {
		l.prev.next = l.next
	}
}

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

var lru *LRUCache

var filename string

// Counter is safe to use concurrently.
type Counter struct {
	mux sync.Mutex
}

var lruCapacity = 1000

// Make a new KeyValue type that is a typed collection of fields
// (Key and Value), both of which are of type string
type KeyValue struct {
	Key, Value, TimeStamp string
}

type KeyValuePair struct {
	Key, Value string
}

// This should be fetched from file (We need a persistent store)
var keyValueStore []KeyValue

// GetToDo takes a string type and returns a ToDo
func (c *Counter) GetKey(key string, value *string) error {

	// use cache (may be later)
	// find the key in file and return
	// return null if not found
	fmt.Printf("calling getKey\n")

	val := lru.Get(key)
	if val != "" {
		*value = val;
		fmt.Printf("cache hit\n")
		return nil;
	}

	fmt.Printf("cache miss\n")

	file, err := os.Open(config[serverIndex]["filename"])
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

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
			lru.Put(key, *value)
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//PutKey ...  TODO: can change timestamp type to Time instead of string
func (c *Counter) PutKey(keyValue KeyValuePair, oldValue *string) error {

	c.mux.Lock()
	defer c.mux.Unlock()
	//get file and find the given key
	//initialise the epoch timestamp
	// first line of the file should contain last updated timestamp
	// update last updated timestamp and the key timestamp
	// if EOF reached append new key, value and timestamp in the end of the file
	// send async requests to update the key value pair with timestamp in other replicas (if no response received
	// in callback from the other server, reinit the server that is not up)
	// fmt.Printf("calling putKey\n")
	keyFound := false
	filePath := config[serverIndex]["filename"]
	curTimeStamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	newKeyValueString := string(keyValue.Key + "," + keyValue.Value + "," + curTimeStamp)

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
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

	lru.Put(keyValue.Key, keyValue.Value)

	for index := range config {
		if index != serverIndex {
			client, err := rpc.DialHTTP("tcp", config[index]["host"]+":"+config[index]["port"])
			if err != nil {
				// callback (check reply and restart server if there is a connection error)
				fmt.Printf("%s ", err)
				RestartServer(index)
			} else {
				defer client.Close()
				client.Go("Counter.SyncKey", KeyValue{Key: keyValue.Key, Value: keyValue.Value, TimeStamp: curTimeStamp}, nil, nil)
			}
		}
	}

	return nil
}

func RestartServer(serverIndex int) {
	// here -r is for server restart
	// cmd := exec.Command("./server", strconv.Itoa(serverIndex), " -r", " &")
	// err := cmd.Start()
	// if err != nil {
	// 	fmt.Printf("error\n")
	// 	fmt.Println(err)
	// }
	// pid := cmd.Process.Pid
	// fmt.Printf("Server %d restarts with process id: %d\n", serverIndex, pid)
}

func SyncReplicas(time int64) error {
	// get all the key value and timestamp pairs from other servers that were updated after this timestamp
	// call synckey method for all the key value pairs

	var updates []KeyValue

	//TODO: make a connection and an RPC Call here? (also need to get updates from all servers)
	for index := range config {
		if index != serverIndex {
			client, err := rpc.DialHTTP("tcp", config[index]["host"]+":"+config[index]["port"])
			if err != nil {
				// RestartServer(index)
				//log.Fatal(err)
			} else {
				defer client.Close()
				// fmt.Println("calling get Updates ++++++++++++++++++++++ ")
				err := client.Call("Counter.GetUpdates", time, &updates)
				if err != nil {
					fmt.Println("error in syncReplica", err)
				} else {
					// To Do - add functionality for bulk update
					// fmt.Println("updates length ==> ", len(updates))
					for _, update := range updates {
						SyncKeyLocally(update)
					}
				}
			}
		}
	}

	return nil
}

//GetUpdates ... to get key values updated after the timestamp
func (c *Counter) GetUpdates(timestamp int64, updates *[]KeyValue) error {
	// return an array of all key,value and timestamp where timestamp > given timestamp
	// fmt.Println("get updates timestamp ===> ", timestamp)

	file, err := os.Open(config[serverIndex]["filename"])
	defer file.Close()
	if err != nil {
		fmt.Println("err 1", err)
		return err
	}

	scanner := bufio.NewScanner(file)

	// Ignore the first line as it contains the last updated time

	if scanner.Scan() {
		scanner.Text()
	}

	for scanner.Scan() {
		var line []string = strings.Split(scanner.Text(), ",")
		if len(line) == 3 {
			TimestampInFile, _ := strconv.ParseInt(line[2], 10, 64)
			if TimestampInFile >= timestamp {
				*updates = append(*updates, KeyValue{Key: line[0], Value: line[1], TimeStamp: line[2]})
			}
		}
	}

	// fmt.Println("updates in getupdates ===> ", len(*updates))

	if err := scanner.Err(); err != nil {
		fmt.Println("err 2", err)
		return err
	}

	return nil
}

func SyncKeyLocally(keyValue KeyValue) error {
	newKeyValueString := keyValue.Key + "," + keyValue.Value + "," + keyValue.TimeStamp

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}

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

	if !found {
		lines = append(lines, newKeyValueString)
	}

	TimeInFile, err = strconv.ParseInt(lines[0], 10, 64)

	if TimeInFile < UpdatedTime {
		lines[0] = strconv.FormatInt(UpdatedTime, 10)
	}

	// fmt.Println()

	// fmt.Println("Num Lines ", strconv.Itoa(len(lines)))

	newFileContent := strings.Join(lines[:], "\n")
	err = ioutil.WriteFile(filename, []byte(newFileContent), 0)

	if err != nil {
		fmt.Printf("%s ", err)
		return err
	}

	return nil
}

func (c *Counter) SyncKey(keyValue KeyValue, reply *string) error {
	// find key in the file
	// if found: check the updated timestamp for the key
	// if updated timestamp > timestamp arg (do nothing)
	// if equal ? (CASE NEEDS TO BE GIVEN A THOUGHT)
	// else: update the key with the passed value and timestamp

	// fmt.Println("calling synckey")
	c.mux.Lock()
	defer c.mux.Unlock()
	SyncKeyLocally(keyValue)
	return nil
}

//Init ... takes in config and index of the current server in config
func Init(index int, restart bool) error {
	// create file and add first line if not already present
	// sync after all servers are up

	a := new(LRUCache)
	a.Init(lruCapacity)
	lru = a

	c := Counter{}
	// Publish the receivers methods
	err := rpc.Register(&c)
	if err != nil {
		fmt.Println("Format of service Task isn't correct. ", err)
	}

	// Sync with the other server before restart
	var time int64 = 0
	if restart == true {
		// read the last updated timestamp (-5 seconds etc -- for failover delay) from the file
		file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		var timestamp string = scanner.Text()
		time, err = strconv.ParseInt(timestamp, 10, 64)
		time = time - 20*1000*60*1000*1000 // 20 minutes
	}

	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234

	listener, e := net.Listen("tcp", config[index]["host"]+":"+config[index]["port"])
	if e != nil {
		fmt.Println("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", config[index]["port"])

	//fmt.Println("restart ===> ", restart)

	if restart == true {
		// fmt.Println("calling sync replica")
		// fmt.Println("Calling Sync replicas =======================")
		err = SyncReplicas(time)
	}

	// Start accept incoming HTTP connections
	// fmt.Println("Listening on the port ------------------------ ")
	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println("Error serving: ", err)
	}
	return nil
}

func main() {
	args := os.Args[1:]
	//fmt.Printf("len %d", len(args))
	serverIndex, _ = strconv.Atoi(args[0])
	pid := os.Getpid()
	fmt.Printf("Server %d starts with process id: %d\n", serverIndex, pid)
	filename = config[serverIndex]["filename"]
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// fmt.Println("creating new file", filename)
		_, err := os.OpenFile(filename, os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Failed to create file ", err)
		}
		curTimeStamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		err = ioutil.WriteFile(filename, []byte(curTimeStamp), 0)
	}
	var restart bool = false
	if len(args) > 1 {
		restart = true
	}
	err = Init(serverIndex, restart)
}
