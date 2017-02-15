package main

import (
	"fmt"
	"github.com/ansgarschmidt/gosomatic/gateways"
	"time"
	"os"
)

func main() {
	fmt.Println("Gosomatic")
	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)
	go gateways.Carbon()
	//gateways.Cloudant()
	for {
		time.Sleep(1000 * time.Second)
	}
}
