package main

import (
	"container/list"
//	"time"
	"utils"
	"log"
//	"list"
)

func main() {
	//timer := time.NewTimer(time.Second * 5)
	
	//k := utils.Get()
	k := 1
	log.Println(k)
	set := utils.NewHashSet()
	log.Println(set.Contains(k))
	
	
	l := list.New()
	l.PushBack("one")
	l.PushBack(2)
	
	for it := l.Front(); it != nil;  {
		log.Println("item:", it.Value)	
		tmp := it
		l.Remove(tmp)
		it = l.Front()
	}
	log.Println(l.Len())
	
}
