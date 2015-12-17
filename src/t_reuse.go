package main
import (
	"fmt"
	"sync"
	"time"
	"math/rand"
	"container/list"
//	"funcTest"
)
type wInt struct {
	data int
}

type node struct {
	count int
	* list.List
}

var send chan int = make(chan int)
var ll *sync.Mutex
var count int
var mp map[int]*node
var r  int

func read() {
	for {
		num := <- send
		t := rand.Intn(3)
		//println("read:", num, "duration", t)
		
		addData(num, t + 1)
	}
}
func printTime(t *time.Time) {
	year,mon,day := t.UTC().Date()
    hour,min,sec := t.UTC().Clock()
	zone,_ := t.UTC().Zone()
	fmt.Printf("UTC time is %d-%d-%d %02d:%02d:%02d %s\n",
                year,mon,day,hour,min,sec,zone)
}

func addData(k int, t int) {
	ll.Lock()
	defer ll.Unlock()
	
	for _, v := range mp {
		if v.count < r {
			v.count++
			dst := time.Now()
		//	printTime(&dst)
			dst = dst.Add(time.Duration(1000 * 1000 * 1000 * t))
		//	printTime(&dst)
			v.PushBack(&dst)
			return
		}
	}
	count++
	item := new(node)
	item.count = 1 
	item.List = list.New()
	dst := time.Now()
//	printTime(&dst)
	dst = dst.Add(time.Duration(1000 * 1000 * 1000 * t))
//	printTime(&dst)
	item.PushBack(&dst)
	mp[k] = item
}

func processData() {
	// 遍历复用流，剔除超时的流
	for {
		ll.Lock()
		println("check")
		now := time.Now()
//		printTime(&now)
		for k, v := range mp {
			for e := v.Front(); e != nil; e = e.Next() {
				data := e.Value.(*time.Time)
				if data.Before(now) {
					v.Remove(e)
					v.count--
				}
			}
			if v.count <= 0 {
				delete(mp, k)
			}
		}

		ll.Unlock()
		time.Sleep(1000 * 1000 * 10)
		println("count", count, "map", len(mp))
	}
}
func initVuser(k int, n int) {
	for i := 0; i < n; i++ {
		p := rand.Intn(10)
		d := time.Duration(1000 * 1000 * 1000 * (p + 1))
		println(k, i, p, d.Seconds())
		time.Sleep(d)
		send <- k * n + i
//		send <- p
	}
}
func main() {
	ll = new(sync.Mutex) 
	mp = map[int]*node{}
	n := 20	//每个虚拟用户的连续请求数
	c := 100	//并发度参数
	r = 8
	for i := 0; i < c; i++ {
		go initVuser(i, n)
	}
	go read()
	processData()
}
