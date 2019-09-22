package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

var serverCount int = 3;
func main() {
	var path string = "./server/server.go";
	for index := 0; index < serverCount; index++ {
		cmd := exec.Command("go", "run", path, strconv.Itoa(index), " &")
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			fmt.Printf("error\n")
			log.Fatal(err)
		}
		pid := cmd.Process.Pid
		fmt.Printf("Server %d starts with process id: %d\n", index, pid)
	}
}
