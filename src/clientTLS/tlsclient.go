package main

import (
        "crypto/tls"
	"fmt"
	"spdy"
//	"io"
	"net/http"
	"os"
	"io/ioutil"
	"log"
)

var spdyClient *spdy.Client

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := spdyClient.Do(r)
	defer resp.Body.Close()
	handle(err)
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	log.Println(resp.StatusCode, resp.ContentLength)
	w.WriteHeader(resp.StatusCode)
	result, err := ioutil.ReadAll(resp.Body)
	handle(err)
	log.Println("here", len(result))
	w.Write(result)
}
func main() {
        cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
        if err != nil {
                fmt.Printf("server: loadkeys: %s", err)
        }
        config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true, NextProtos: []string{"spdy/3"}}
        conn, err := tls.Dial("tcp", "127.0.0.1:4040", &config)
        if err != nil {
                fmt.Printf("client: dial: %s", err)
        }
	client, err := spdy.NewClientConn(conn)
	handle(err)
	spdyClient = client
/*	req, err := http.NewRequest("GET", "http://www.sohu.com/banana", nil)
	handle(err)
	res, err := client.Do(req)
	handle(err)
	data := make([]byte, int(res.ContentLength))
	_, err = res.Body.(io.Reader).Read(data)
	fmt.Println(string(data))
	res.Body.Close()*/

	http.HandleFunc("/", handler)
	log.Println("Start serving on port 1234")
	http.ListenAndServe(":1234", nil)
	os.Exit(0)
}
