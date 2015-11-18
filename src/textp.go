package main
import ( 
	"strings"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func checkStatic(oname string) bool {
	if strings.HasSuffix(oname, ".png") {
		return true
	} else if strings.HasSuffix(oname, ".jpg") {
		return true
	} else if strings.HasSuffix(oname, ".ico") {
		return true
	} else if strings.HasSuffix(oname, ".bmp") {
		return true
	} else if strings.HasSuffix(oname, ".flv") {
		return true
	} else if strings.HasSuffix(oname, ".swf") {
		return true
	} else if strings.HasSuffix(oname, ".css") {
		return true
	} else if strings.HasSuffix(oname, ".js") {
		return true
	} else if strings.HasSuffix(oname, ".gif") {
		return true 
	} else {
		return false
	}
}

func getOnameFromUrl(url string) string {
	items := strings.Split(url, "/")
	//println(url, items[len(items) - 1])
	return items[len(items) - 1]
}

func match(urls []string, filepaths []string, filedatas []string) {
	println(len(filepaths), len(filedatas))
	all := 0
	countF := 0
	for i, url := range urls {
		oname := getOnameFromUrl(url)

		if checkStatic(oname) == true {
			all++
			for j, data := range filedatas {
				if strings.Count(data, oname) > 0 {
					countF++
					println(i, url, filepaths[j])
					break
				}
				if j == len(filedatas) - 1 {
					println("not found:", url)
				}
			}
		}
	}
	println("all:", all, "find:", countF)
}

func readFile(filepath string) []byte {
	fi, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return fd
}

func getUrls(filePath string) []string {
	data := readFile(filePath)
	urls := []string{}
	sdata := string(data)
	count := 0

	for {
		start := strings.Index(sdata[count:], "url\":")
		end := strings.Index(sdata[count + start:], "\",")
		if end < 7 {
			break
		}
		urls = append(urls, sdata[count + start - 1 : count + start + end])
		count = count + start + end
		if count >len(sdata) {
			break
		}
	}
	return urls
}

func getFiles(dirPath string)([]string, []string){
	filenames := []string{}
	filedatas := []string{}
	filepath.Walk(dirPath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			} else {
//				println(f.Size())
			}
			filenames = append(filenames, path)
			//println(path)
			return nil 
		})
	for _, filename := range filenames {
		data := readFile(filename)
		sdata := string(data)
		filedatas = append(filedatas, sdata)
	}
	return filenames, filedatas
}
func main() {
	urls :=  getUrls("/home/ww/cache/data/har/www.163.com.har")
	filepaths, filedatas := getFiles("/home/ww/cache/data/163")
	match(urls, filepaths, filedatas)
}
