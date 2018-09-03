package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestInitClient(t *testing.T) {
	var client *http.Client
	client =  InitClient("socks5://127.0.0.1:1080")
	resp, err := client.Get("http://ac168.info/bt/thread.php?fid=4")
	if err != nil{
		fmt.Println("get 方法错误")
		/*panic(err)*/
	}
	defer resp.Body.Close()
	fmt.Println(resp.Body)
}