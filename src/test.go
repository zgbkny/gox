package main

import (
	//"container/list"
//	"time"
	//"utils"
	"log"
//	"list"
)

func main() {
	//timer := time.NewTimer(time.Second * 5)
	
	//k := utils.Get()
    k := 100
	arr := make([]int, k)
	log.Println(len(arr))
    arra := make([]interface{}, 10)
    log.Println(arra[1])	
    if arra[1] == nil {
        log.Println("is nil")
    } else {
        log.Println("is not nil")
    }
}
