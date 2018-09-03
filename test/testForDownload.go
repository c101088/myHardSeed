package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"net/url"
	"io"
	"bytes"
)

func main(){
	//testUrl:="http://www.jandown.com/fetch.php"
//	testHttpGet(testUrl)
	//testHttpClientGet(testUrl)
	//testHttpPost(testUrl)
	testHttpClientPost("http://www.rmdown.com/download.php?ref="+"182b088e2c6877068eed7b85e00a43c90783d6362c3"+"&reff="+"106100"+"&submit=download")
//	testGetPicture("http://imglink.ru/pictures/23-08-18/1bb9c3118bb3d118ebd973a5d61130e0.jpg")
}


func writeToFile(filename string,content string ){
	f,_:= os.Create(filename) //创建文件
	defer f.Close()
	f.WriteString(content) //写入文件(字节数组)
	f.Sync()
}

func testGetPicture(testUrl string)  {
	urli := url.URL{}
	urlProxy,_:= urli.Parse("socks5://127.0.0.1:1080")
	client := &http.Client{
		Transport:&http.Transport{Proxy:http.ProxyURL(urlProxy)},
	}
	resp, err := client.Get(testUrl)
	if err !=nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
	fmt.Printf("%v\n",body)
	resp1 :=string(body)
	var f *os.File
	f, _ = os.Create("D:/360Combo/xxx.jpg")

	if f != nil && len(resp1)>0{
		io.Copy(f,bytes.NewReader([]byte(resp1)))
	}
	defer f.Close()

}

/*func testNewRequest(testUrl string){
	urli := url.URL{}
	urlProxy,_:= urli.Parse("socks5://127.0.0.1:1080")
	client := &http.Client{
		Transport:&http.Transport{Proxy:http.ProxyURL(urlProxy)},
	}
	req, e := http.NewRequest("POST",)
}*/

func testHttpClientPost(testUrl string){
	urli := url.URL{}
	urlProxy,_:= urli.Parse("socks5://127.0.0.1:1080")
	client := &http.Client{
		Transport:&http.Transport{Proxy:http.ProxyURL(urlProxy)},
	}
	resp, err := client.Get(testUrl)
	//	resp, err := client.Post(testUrl,"application/x-www-form-urlencoded	",strings.NewReader("ref=182b088e2c6877068eed7b85e00a43c90783d6362c3"))

	if err!= nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
	fmt.Printf("%s\n",body)
	resp1 :=string(body)
	var f *os.File
	f, _ = os.Create("D:/360Combo/xxx.torrent")

	if f != nil && len(resp1)>0{
		io.Copy(f,bytes.NewReader([]byte(resp1)))
	}
	defer f.Close()

}

func testHttpPost(testUrl string) {

		resp, err := http.Post(testUrl,"application/x-www-form-urlencoded",strings.NewReader("name=abc"))
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		}
		writeToFile("testHttpPost.html",string(body))
}

func testHttpClientGet(testUrl string) {

	urli := url.URL{}
	urlProxy,_:= urli.Parse("socks5://127.0.0.1:1080")
	client := &http.Client{
		Transport:&http.Transport{Proxy:http.ProxyURL(urlProxy)},
	}
	resp, err := client.Get(testUrl)
	if err!= nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
	fmt.Println(string(body))
}

func testHttpGet(testUrl string){
	resp, err := http.Get(testUrl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
	fmt.Println(string(body))
}
