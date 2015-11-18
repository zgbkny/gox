
package main
import "fmt"
//import "tcpServer"
//import "funcTest"
import "utils"
import (
//	"container/list"
)
func write1() {


	for i:= 0; i < 100; i++ {
		return
	}
	
}

func write2() {

	for i:= 0; i < 100; i++ {
		
		return
	}
	

}

func read(data string) {
	fmt.Println("read:", data)
}
func main() {
	fmt.Println("hello, world!")
//	funcTest.FuncTest(read)
	var i uint32
	i = 578
	data := utils.Int32ToBytes(i)
	k := 23845
	dd := utils.Int16ToBytes(k)
	fmt.Println(dd)
	fmt.Println(utils.BytesToInt16(dd))
	t := utils.BytesToInt32(data[2:4])
	fmt.Println(i)
	fmt.Println(data[2 : 4])
	fmt.Println(t)

}
