package main

import "C"

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"
	"unsafe"
)

var err error
var client *rpc.Client
var serverIndex int
var reply string
var serverList []string

//KeyValuePair ... interface type
type KeyValuePair struct {
	Key, Value string
}

//export kv739_init
func kv739_init(cserverListArg **C.char, length C.int) C.int {
	//TODO: can you work without length argument?

	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(cserverListArg))[:length:length]
	serverListArg := make([]string, length)
	for i, s := range tmpslice {
		serverListArg[i] = C.GoString(s)
	}
	fmt.Printf("serverListArg %s\n", serverListArg[0])
	rand.Seed(time.Now().UnixNano())
	serverIndex = rand.Intn(len(serverListArg))
	fmt.Printf("server index ====> %d\n", serverIndex)
	serverList = serverListArg
	address := serverList[serverIndex]
	client, err = rpc.DialHTTP("tcp", address)
	if err != nil {
		fmt.Println("Connection error: ", err)
		return C.int(-1)
	}
	return C.int(0)
}

//export kv739_shutdown
func kv739_shutdown() C.int {
	err := client.Close()
	if err != nil {
		fmt.Println("Unable to shutdown client connection: ", err)
		return C.int(-1)
	}
	return C.int(0)
}

//export kv739_get
func kv739_get(ckey *C.char, cvalue *C.char) C.int {
	key := C.GoString(ckey)
	value := C.GoString(cvalue)
	err := client.Call("Task.GetKey", key, &value)
	if err != nil {
		fmt.Println("Could not get key: ", key, err)
		//Retry logic
		if len(serverList) > 1 {
			for index, server := range serverList {
				if index != serverIndex {
					client, err = rpc.DialHTTP("tcp", server)
					if err == nil {
						err := client.Call("Task.GetKey", key, &value)
						if err == nil {
							serverIndex = index
							if len(value) > 1 {
								convertGoToString(cvalue, value)
								return C.int(0)
							}
							return C.int(1)
						}
						fmt.Println("Unable to get key: ", key, " from server: ", server, " err: ", err)
					} else {
						fmt.Println("Unable to establish connection with server: ", server, err)
					}
				}
			}
		}
		return C.int(-1)
	}
	if len(value) > 1 {
		convertGoToString(cvalue, value)
		return C.int(0)
	}
	return C.int(1)
}

//export kv739_put
func kv739_put(ckey *C.char, cvalue *C.char, coldValue *C.char) C.int {
	key := C.GoString(ckey)
	value := C.GoString(cvalue)
	oldValue := C.GoString(coldValue)
	err := client.Call("Task.PutKey", KeyValuePair{Key: key, Value: value}, &oldValue)
	if err != nil {
		fmt.Println("Could not put key: ", key, " value: ", value, " err: ", err)
		//TODO: Retry logic only if err contains connection
		if len(serverList) > 1 {
			for index, server := range serverList {
				if index != serverIndex {
					client, err = rpc.DialHTTP("tcp", server)
					if err == nil {
						err := client.Call("Task.PutKey", KeyValuePair{Key: key, Value: value}, &oldValue)
						if err == nil {
							serverIndex = index
							if len(oldValue) > 1 {
								convertGoToString(coldValue, oldValue)
								return C.int(0)
							}
							return C.int(1)
						}
						fmt.Println("Unable to put key: ", key, " on server: ", server, " err: ", err)
					} else {
						fmt.Println("Unable to establish connection with server: ", server, err)
					}
				}
			}
		}
		return C.int(-1)
	}

	if len(oldValue) > 1 {
		convertGoToString(coldValue, oldValue)
		return C.int(0)
	}
	return C.int(1)
}

func convertGoToString(coldValue *C.char, oldValue string) {
	lenvalue := len(oldValue)
	gData := (*[1 << 30]byte)(unsafe.Pointer(coldValue))[:lenvalue:lenvalue]
	copy(gData[0:], oldValue)
}

func main() {
}
