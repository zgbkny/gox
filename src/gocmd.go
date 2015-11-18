package main
import (
	"fmt"
	"os/exec"
	"strconv"
)

func genFixSizeFile(inf string, outf string, size int) {
	fmt.Println(inf, outf, size)
	cmd := exec.Command("/bin/bash", "-c", "dd if=" + inf + " of=" + outf + " bs=" + strconv.Itoa(size) + " count=1")
	//cmd := exec.Command("dd", "if=", "test.go", "of=", "test", "bs=", "24", "count =", "1")
	d, err := cmd.Output()
	fmt.Println(err, string(d))
}
func genDir(dir string) {
	cmd := exec.Command("mkdir", dir)
	d, _ := cmd.Output()
	fmt.Println(string(d))
}

/*产生特定格式的文件*/
func genFiles() {
	inFile := "/home/ww/cache/im*k.tar.gz"
	outFile := "/home/ww/cache/"
	for i := 1; i <= 20; i++ {
		foldName := strconv.Itoa(i * 20)
		flatName := new([5]byte)
		flatName[0] = 'i'
		fmt.Println(foldName, foldName[0], len(foldName))
		if len(foldName) < 3 {
			flatName[1] = '0'
			flatName[2] = foldName[0]
			flatName[3] = foldName[1]
		} else {
			flatName[1] = foldName[0]
			flatName[2] = foldName[1]
			flatName[3] = foldName[2]
		}
		flatName[4] = 'k'
		fmt.Println(string(flatName[:5]))
		genDir(outFile + string(flatName[:5]))
		for j := 0; j < 10; j++ {
			
			genFixSizeFile(inFile, outFile + string(flatName[:5]) + "/test00_000" + strconv.Itoa(j) + ".png", i * 20 * 1024)
		}
	}
}


func main() {
	genFiles()
}	
