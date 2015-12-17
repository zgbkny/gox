package main

import (
//	"container/list"
//	"time"
	"utils"
	"log"
)

func main() {
	//timer := time.NewTimer(time.Second * 5)
	
	//k := utils.Get()
	k := 1
	log.Println(k)
	set := utils.NewHashSet()
	log.Println(set.Contains(k))
	
}
