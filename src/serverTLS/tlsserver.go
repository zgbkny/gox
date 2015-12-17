package main

import (
	"fmt"
	"spdy"
	"net/http"
//	"os"
	"io/ioutil"
	"log"
	//"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = "http"
	r.URL.Host = "www.sohu.com"

	resp, err := http.DefaultClient.Do(r)
	defer resp.Body.Close()
	if err != nil { panic(err) }
	
	for k, v := range resp.Header {
		for _, vv := range v {
			log.Println(k, vv)
		}
	}
	
	log.Println(resp.StatusCode, resp.ContentLength)
	w.WriteHeader(resp.StatusCode)
	result, err := ioutil.ReadAll(resp.Body)
//	if err != nil { panic(err) }
	w.Write(result)
	log.Println("here 1", len(result))
}

func main() {
 //       spdy.EnableDebug()
        http.HandleFunc("/", handler)
	err := spdy.ListenAndServeTLS("localhost:4040", "server.pem", "server.key" , nil)
	if err != nil {
		fmt.Println(err)
	}
}


