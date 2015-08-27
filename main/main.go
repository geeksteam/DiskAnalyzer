package main

import (
	"log"
	"time"
)

const testDir = "/home/max/Documents/fld"

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds)
	for i := 0; i < 10; i++ {
		go log.Println("hello")
	}
	time.Sleep(3000)
}
