package main
import "fmt"
import "sync"
//import "tcpServer"



var send chan []byte = make(chan []byte)

var ll sync.Mutex 

var k int = 0
var j int = 0

func send1(send chan []byte, data []byte) {
	//ll.Lock()
	//defer ll.Unlock()
	send <- data[:]	
}
func write1() {
	var data [4]byte
	data[0] = 1
	data[1] = 2
	data[2] = 1
	data[3] = 2

	for i:= 0; i < 100; i++ {
		send1(send, data[:])
		k = i
	}
}

func write2() {
	var data [4]byte
	data[0] = 4
	data[1] = 5
	data[2] = 4
	data[3] = 5
	for i:= 0; i < 100; i++ {
		send1(send, data[:])	
		j = i
	}
	

}

func read() {
	for i:= 0; i < 200; i++ {
		buf, ok:= <- send
		if !ok {
			break
		}
		fmt.Println("hello, world!" , buf, i, k, j)
	}
}

func main() {
	fmt.Println("hello, world!")
	
	go write1()
	go write2()
	read()
}
