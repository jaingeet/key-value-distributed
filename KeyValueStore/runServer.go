package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

var serverCount int = 3

func main() {
<<<<<<< HEAD
	var path string = "./server/server.go"
	for index := 1; index <= serverCount; index++ {
		cmd := exec.Command("go", "run", path, strconv.Itoa(index), " &")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
=======
	var path string = "./server/server.go";
	for index := 0; index < serverCount; index++ {
		cmd := exec.Command("go", "run", path, strconv.Itoa(index), " &")
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		err := cmd.Start()
>>>>>>> 124a61f46453650b29d7a4ec27467884504cce1b
		if err != nil {
			fmt.Printf("error\n")
			log.Fatal(err)
		}
		pid := cmd.Process.Pid
		fmt.Printf("Server %d starts with process id: %d\n", index, pid)
	}
}
