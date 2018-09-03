package parserWebpage

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"self-goProject/myHardSeed/cmdInteracion"
	"strconv"
	"strings"
)


//todo 处理aicheng的两层解析和caoliu的三层解析
type AichengWebpage struct {
	cmdInteracion.CmdOrder
	AiSeedAndPicChan chan SeedAndPicChan
	ItemNum int
}

type CaoliuWebpage struct {
	cmdInteracion.CmdOrder
	AiSeedAndPicChan chan SeedAndPicChan
	ItemNum int
}


type SeedAndPicChan struct {
	SeedCode string
	Name string
	PicSrc []string
}

const ichengPageRef = `<h3><a href="(htm_data[^.]+.html)"[^>]*>([^<]*)</a></h3>`
//<img src="http://imagizer.imageshack.com/img921/9979/jGt1Dq.jpg" border="0"
const ichengPicRef =`<img src="([^"]+)"`
//>http://www.jandown.com/link.php?ref=CoSK9BaN</a>
const ichengSeedRef =`>http[^\?]+\?ref=([^<]+)<`
//<h3><a href="htm_data/15/1808/3233754.html" target="_blank" id="">[FHD/5.07G]IPZ-662 很火的居酒屋店員 希島あいり</a></h3>
const caoliuPageRef =`<h3><a href="(htm_data[^.]+.html)"[^>]*>([^<]*)</a></h3>`
//<img data-link='http://img599.net/image/wAg8' data-src='http://img599.net/images/2013/05/29/snis00580jp-1.th.jpg'>&nbsp;
const caoliuPicRef = `data-src='([^']+)'>`
//>http://www.rmdown.com/link.php?hash=182cd6b21b984e42e6d82e44a99e2a9e4eff200e440</a>
const caoliuSeedRef  =`>http[^\?]+\?hash=([^<]+)<`
//<INPUT TYPE="hidden" NAME="reff" value="106100">
const caoliuSeedReff =`<INPUT TYPE="hidden" NAME="reff" value="([^"]+)">`

//获取范围内单个对象的网址
//判断like和hate菜单
func (aichengWebpage *AichengWebpage)ParserWebpage(targetUrlSub string,aichengHost string,client *http.Client){
	pageNum :=1
	itemNum :=0
	itemStart := true
	targetUrl := aichengHost+targetUrlSub
	likeAndHateflag := true
	fmt.Println("正在解析网址：")
forCode:
	for true {
//		log.Printf("pageNum:%v ,%s\n",pageNum,targetUrl + "&page="+ strconv.Itoa(pageNum))
		resp, err := client.Get(targetUrl + "&page="+ strconv.Itoa(pageNum))
		if err !=nil {
			log.Printf("get error whlie fetch %s \n",targetUrl + "&page="+ strconv.Itoa(pageNum))
		}
		bodyReader := bufio.NewReader(resp.Body)
		defer resp.Body.Close()
		e := determineEncoding(bodyReader)
		utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
		body, err := ioutil.ReadAll(utf8Reader)   //请求数据进行读取
		regexRes:= regexp.MustCompile(ichengPageRef)
		temp :=regexRes.FindAllSubmatch(body,-1)
		itemNum += len(temp)


		fmt.Printf("\r")
		processPercent :=itemNum/(aichengWebpage.TopicRange[1]-aichengWebpage.TopicRange[0])
		if processPercent>1 {processPercent = 1}
		fmt.Printf("%s>  %d % \n",strings.Repeat("=",processPercent*10),processPercent*100)


		if itemStart ==true{
			if itemNum>aichengWebpage.TopicRange[0] {
				if len(temp)>=aichengWebpage.TopicRange[1]{
					for i,v := range temp{
						if (i >= aichengWebpage.TopicRange[0]) && (i<=aichengWebpage.TopicRange[1]){
							v[1]=[]byte(aichengHost+string(v[1]))
							likeAndHateflag = likeAndHateCheck(aichengWebpage.TopicLike, aichengWebpage.TopicHate, string(v[2]))
							if likeAndHateflag{
//								fmt.Printf("send a parserSingleWebpage request\n")
								aichengWebpage.parserSingleWebpage(v,client)

							}
						}
					}
					break forCode
				}else{
					pageStartIndex:=aichengWebpage.TopicRange[0]%len(temp)
					for i,v := range temp{
						if i >= pageStartIndex{
							v[1]=[]byte(aichengHost+string(v[1]))
							likeAndHateflag = likeAndHateCheck(aichengWebpage.TopicLike, aichengWebpage.TopicHate, string(v[2]))
							if likeAndHateflag{
//								fmt.Printf("send a parserSingleWebpage request\n")
								aichengWebpage.parserSingleWebpage(v,client)
							}
						}
					}
				}

				itemStart =false
			}
		}else{
			if itemNum>aichengWebpage.TopicRange[1]{
				pageEndIndex := aichengWebpage.TopicRange[1]%len(temp)
				for i,v := range temp{
					if i<= pageEndIndex{
						v[1]=[]byte(aichengHost+string(v[1]))
						likeAndHateflag = likeAndHateCheck(aichengWebpage.TopicLike, aichengWebpage.TopicHate, string(v[2]))
						if likeAndHateflag{
//							fmt.Printf("send a parserSingleWebpage request\n")
							aichengWebpage.parserSingleWebpage(v,client)
						}
					}
				}
				break forCode
			}
			for _,v := range temp{
				v[1]=[]byte(aichengHost+string(v[1]))
				likeAndHateflag = likeAndHateCheck(aichengWebpage.TopicLike, aichengWebpage.TopicHate, string(v[2]))
				if likeAndHateflag{
//					fmt.Printf("send a parserSingleWebpage request\n")
					aichengWebpage.parserSingleWebpage(v,client)
				}
			}
		}
		pageNum++
	}

}
//解析单个网址获得其图片网址和种子码
func (aichengWebpage *AichengWebpage) parserSingleWebpage(value [][]byte,client *http.Client){
	aichengWebpage.ItemNum++
//	fmt.Printf("get a parserSingleWebpage request\n")
	resp, err := client.Get(string(value[1]))
	if err !=nil {
		panic(err)
	}
	bodyReader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	body, err := ioutil.ReadAll(utf8Reader)   //请求数据进行读取
	regexPicRes:= regexp.MustCompile(ichengPicRef)
	tempPic :=regexPicRes.FindAllSubmatch(body,-1)
	regexSeedRes:= regexp.MustCompile(ichengSeedRef)
	tempSeed :=regexSeedRes.FindAllSubmatch(body,-1)
//	fmt.Printf("Name : %s,Item Url :%s \n",value[2],value[1])

	AiSeedAndPicChan :=SeedAndPicChan{}
	AiSeedAndPicChan.Name = string(value[2])
	for _,v := range tempPic{
		AiSeedAndPicChan.PicSrc = append(AiSeedAndPicChan.PicSrc,string(v[1]))
//		fmt.Printf("picture url : %s \n",v[1])
	}
	for _,v := range tempSeed  {
		AiSeedAndPicChan.SeedCode = string(v[1])
//		fmt.Printf("Seed code : %s \n",v[1])
	}

	aichengWebpage.AiSeedAndPicChan<-AiSeedAndPicChan


//	fmt.Printf("send one AiSeedAndPicChan")
}

func (aichengWebpage *AichengWebpage)DownloadPicAndSeed(client *http.Client,out chan int){
//	log.Printf("save path :%s",aichengWebpage.SavePath)
	go func() {
		for  {
			aii:=<-aichengWebpage.AiSeedAndPicChan
			name :=aii.Name
			picSrcUrl := aii.PicSrc
			seedCode := aii.SeedCode

//			fmt.Printf("name : %s, seedCode : %s",name ,seedCode)
			for i:= 0;i<=len(picSrcUrl)-1;i++{
				aidownPic(name+strconv.Itoa(i),aichengWebpage.SavePath,picSrcUrl[i],client)
			}
			aidownSeed(name,aichengWebpage.SavePath,seedCode,client,out)
		}
	}()
}

func aidownSeed(name string ,path string ,code string ,client *http.Client,out chan int){
	resp, err := client.Post("http://www.jandown.com/fetch.php","application/x-www-form-urlencoded	",strings.NewReader("code="+code))

	if err!= nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
	//fmt.Printf("%v\n",body)
//	resp1 :=string(body)
	var f *os.File
	f, _ = os.Create(path + name + ".torrent")
//	fmt.Printf("path: %s\n",path + name + ".torrent")
	if f != nil && len(body)>0{
		io.Copy(f,bytes.NewReader(body))
	}
	defer f.Close()
	out<-1
}

func aidownPic(name string,path string ,url string,client *http.Client){
	resp, err := client.Get(url)
	if err !=nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
//	fmt.Printf("%v\n",body)
	var f *os.File
	f, _ = os.Create(path + name + ".jpg")

	if f != nil && len(body)>0{
		io.Copy(f,bytes.NewReader(body))
	}
	defer f.Close()

}


func (caoliuWebpage *CaoliuWebpage)ParserWebpage(targetUrlSub string,caoliuHost string,client *http.Client){
	pageNum :=1
	itemNum :=0
	itemStart := true
	targetUrl := caoliuHost+targetUrlSub
	likeAndHateflag := true
	fmt.Println("正在解析网址：")
forCode:
	for true {
				log.Printf("%s \n",targetUrl + "&search=&page="+ strconv.Itoa(pageNum))
		resp, err := client.Get(targetUrl + "&search=&page="+ strconv.Itoa(pageNum))
		if err !=nil {
			log.Printf("get error whlie fetch %s \n",targetUrl + "&search=&page="+ strconv.Itoa(pageNum))
		}
		bodyReader := bufio.NewReader(resp.Body)
		defer resp.Body.Close()
		e := determineEncoding(bodyReader)
		utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
		body, err := ioutil.ReadAll(utf8Reader)   //请求数据进行读取
		regexRes:= regexp.MustCompile(caoliuPageRef)
		temp :=regexRes.FindAllSubmatch(body,-1)
		itemNum += len(temp)


		fmt.Printf("\r")
		processPercent :=itemNum/(caoliuWebpage.TopicRange[1]-caoliuWebpage.TopicRange[0])
		if processPercent>1 {processPercent = 1}
		fmt.Printf("%s>  %d % \n",strings.Repeat("=",processPercent*10),processPercent*100)


		if itemStart ==true{
			if itemNum>caoliuWebpage.TopicRange[0] {
				if len(temp)>=caoliuWebpage.TopicRange[1]{
					for i,v := range temp{
						if (i >= caoliuWebpage.TopicRange[0]) && (i<=caoliuWebpage.TopicRange[1]){
							v[1]=[]byte(caoliuHost+string(v[1]))
							likeAndHateflag = likeAndHateCheck(caoliuWebpage.TopicLike, caoliuWebpage.TopicHate, string(v[2]))
							if likeAndHateflag{
								//								fmt.Printf("send a parserSingleWebpage request\n")
								caoliuWebpage.parserSingleWebpage(v,client)

							}
						}
					}
					break forCode
				}else{
					pageStartIndex:=caoliuWebpage.TopicRange[0]%len(temp)
					for i,v := range temp{
						if i >= pageStartIndex{
							v[1]=[]byte(caoliuHost+string(v[1]))
							likeAndHateflag = likeAndHateCheck(caoliuWebpage.TopicLike, caoliuWebpage.TopicHate, string(v[2]))
							if likeAndHateflag{
								//								fmt.Printf("send a parserSingleWebpage request\n")
								caoliuWebpage.parserSingleWebpage(v,client)
							}
						}
					}
				}

				itemStart =false
			}
		}else{
			if itemNum>caoliuWebpage.TopicRange[1]{
				pageEndIndex := caoliuWebpage.TopicRange[1]%len(temp)
				for i,v := range temp{
					if i<= pageEndIndex{
						v[1]=[]byte(caoliuHost+string(v[1]))
						likeAndHateflag = likeAndHateCheck(caoliuWebpage.TopicLike, caoliuWebpage.TopicHate, string(v[2]))
						if likeAndHateflag{
							//							fmt.Printf("send a parserSingleWebpage request\n")
							caoliuWebpage.parserSingleWebpage(v,client)
						}
					}
				}
				break forCode
			}
			for _,v := range temp{
				v[1]=[]byte(caoliuHost+string(v[1]))
				likeAndHateflag = likeAndHateCheck(caoliuWebpage.TopicLike, caoliuWebpage.TopicHate, string(v[2]))
				if likeAndHateflag{
					//					fmt.Printf("send a parserSingleWebpage request\n")
					caoliuWebpage.parserSingleWebpage(v,client)
				}
			}
		}
		pageNum++
	}

}

func (caoliuWebpage *CaoliuWebpage)DownloadPicAndSeed(client *http.Client,out chan int){
	//	log.Printf("save path :%s",aichengWebpage.SavePath)
	go func() {
		for  {
			aii:=<-caoliuWebpage.AiSeedAndPicChan
			name :=aii.Name
			picSrcUrl := aii.PicSrc
			seedCode := aii.SeedCode

			//			fmt.Printf("name : %s, seedCode : %s",name ,seedCode)
			for i:= 0;i<=len(picSrcUrl)-1;i++{
				caodownPic(name+strconv.Itoa(i),caoliuWebpage.SavePath,picSrcUrl[i],client)
			}
			caodownSeed(name,caoliuWebpage.SavePath,seedCode,client,out)
		}
	}()
}



func caodownSeed(name string ,path string ,code string ,client *http.Client,out chan int){



	resp, err := client.Get("http://rmdown.com/link.php?hash="+code)
	if err !=nil {
		log.Printf("get error whlie fetch %s \n","http://rmdown.com/link.php?hash="+code)
	}
	bodyReader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	body, err := ioutil.ReadAll(utf8Reader)   //请求数据进行读取
	regexRes:= regexp.MustCompile(caoliuSeedReff)
	reffTemp :=regexRes.FindAllSubmatch(body,-1)
	var reffStr string
	for _,v := range reffTemp{
		reffStr = string(v[1])
	}
	resp.Body.Close()
	resp, err = client.Get("http://www.rmdown.com/download.php?ref="+code+"&reff="+reffStr+"&submit=download")
	if err!= nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
	//fmt.Printf("%v\n",body)
	//	resp1 :=string(body)
	var f *os.File
	f, _ = os.Create(path + name + ".torrent")
	//	fmt.Printf("path: %s\n",path + name + ".torrent")
	if f != nil && len(body)>0{
		io.Copy(f,bytes.NewReader(body))
	}
	defer f.Close()
	out<-1
}

func caodownPic(name string,path string ,url string,client *http.Client){
	resp, err := client.Get(url)
	if err !=nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)   //请求数据进行读取
	if err != nil {
		// handle error
	}
	//	fmt.Printf("%v\n",body)
	var f *os.File
	f, _ = os.Create(path + name + ".jpg")

	if f != nil && len(body)>0{
		io.Copy(f,bytes.NewReader(body))
	}
	defer f.Close()

}




//解析单个网址获得其图片网址和种子码
func (caoliuWebpage *CaoliuWebpage) parserSingleWebpage(value [][]byte,client *http.Client){
	caoliuWebpage.ItemNum++
	//	fmt.Printf("get a parserSingleWebpage request\n")
	resp, err := client.Get(string(value[1]))
	if err !=nil {
		panic(err)
	}
	bodyReader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	body, err := ioutil.ReadAll(utf8Reader)   //请求数据进行读取
	regexPicRes:= regexp.MustCompile(caoliuPicRef)
	tempPic :=regexPicRes.FindAllSubmatch(body,-1)
	regexSeedRes:= regexp.MustCompile(caoliuSeedRef)
	tempSeed :=regexSeedRes.FindAllSubmatch(body,-1)
	//	fmt.Printf("Name : %s,Item Url :%s \n",value[2],value[1])

	AiSeedAndPicChan :=SeedAndPicChan{}
	AiSeedAndPicChan.Name = string(value[2])
	for i,v := range tempPic{
		if i>5 {break}
		AiSeedAndPicChan.PicSrc = append(AiSeedAndPicChan.PicSrc,string(v[1]))
		//		fmt.Printf("picture url : %s \n",v[1])
	}
	for _,v := range tempSeed  {
		AiSeedAndPicChan.SeedCode = string(v[1])
		//		fmt.Printf("Seed code : %s \n",v[1])
	}

	caoliuWebpage.AiSeedAndPicChan<-AiSeedAndPicChan


	//	fmt.Printf("send one AiSeedAndPicChan")
}
func likeAndHateCheck(likeList []string,hateList []string,seedName string) bool  {
	result := true
	if len(likeList) == 0{
		result = true
		for _,v:= range hateList{
			b := strings.Contains(seedName, v)
			if b==true{
				result =false
			}
		}
	}else{
		result = false
		for _,v:=range likeList{
			b := strings.Contains(seedName, v)
			if b ==true{
				result = true
			}
		}
		for _,v := range hateList{
			b := strings.Contains(seedName, v)
			if b == true {
				result = false
			}
		}
	}
/*	if result == true {
		log.Printf("congratulations! you got a like seed!\n")
	}*/
	return result
}
func ForTestingRegx(){
	urli := url.URL{}
	urlProxy,_:= urli.Parse("socks5://127.0.0.1:1080")
	client := &http.Client{
		Transport:&http.Transport{Proxy:http.ProxyURL(urlProxy)},
	}
	resp, err := client.Get("http://ac168.info/bt/htm_data/4/1808/964515.html")
	if err !=nil {
		log.Printf("get error when fetch url :  %s \n","http://ac168.info/bt/htm_data/16/1808/964515.html")
	}
	defer resp.Body.Close()
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	body, err := ioutil.ReadAll(utf8Reader)   //请求数据进行读取
//	fmt.Printf("%s",body)
	regexPicRes:= regexp.MustCompile(ichengPicRef)
	tempPic :=regexPicRes.FindAllSubmatch(body,-1)
	for _,v := range tempPic{
		fmt.Printf("%s,picture url : %s \n",v[0],v[1])
	}
	regexSeedRes:= regexp.MustCompile(ichengSeedRef)
	tempSeed :=regexSeedRes.FindAllSubmatch(body,-1)

	for _,v := range tempSeed  {
		fmt.Printf("Seed code : %s \n",v[1])
	}

}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	byte, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(byte, "")
	return e
}
