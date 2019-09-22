package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var serverCount int = 3;
func main() {
	var path string = "./server/server.go";
	for index := 1; index <= serverCount; index++ {
		cmd := exec.Command("go", "run", path, strconv.Itoa(index), "&")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error\n")
			log.Fatal(err)
		}
	}
}
