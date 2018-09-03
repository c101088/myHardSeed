package main

import (
	"fmt"
	"github.com/godfather1103/hardseed-go/utils"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"self-goProject/myHardSeed/cmdInteracion"
	"self-goProject/myHardSeed/parserWebpage"
	"strconv"
	"time"
)

const (
	caoliuHost  = "http://www.t66y.com/"
	aichengHost ="http://www.ac168.info/bt/"
)

func main(){
	var cmdLineOption cmdInteracion.CmdOrder
	err:= cmdInteracion.ParserCmd(&cmdLineOption)
	if err != nil {
		panic(nil)
	}

	//Show the content of cmdOrder
	showCmdOrder(cmdLineOption)
	//check the aid URL
	targetUrlStrSub,aichengFlag,err:=getParserTargetUrl(cmdLineOption.AvClass)
	if err != nil {
		panic(err)
	}
	// initiatize the client with proxy
	var client *http.Client
	client = InitClient(cmdLineOption.ProxyString)
	//set the cpu core for hardseed
	runtime.GOMAXPROCS(cmdLineOption.CoreNum)

	//set the download path of files
	downloadDirName:= cmdLineOption.SavePath+"[" + cmdLineOption.AvClass+strconv.Itoa(cmdLineOption.TopicRange[0]) + "~" + strconv.Itoa(cmdLineOption.TopicRange[1])+ "]"+ func()string {
		if runtime.GOOS == "windows" {
			return "\\"}else{
			return "/"}}()
	//run the main body of hardseed

	seedAndPic := make(chan parserWebpage.SeedAndPicChan,10)
	out := make(chan int)
	if aichengFlag{
		aichengWebpage := parserWebpage.AichengWebpage{cmdLineOption,seedAndPic,0}
		aichengWebpage.SavePath = downloadDirName
//		fmt.Println(aichengWebpage.SavePath)
		exists,_ := utils.PathExists(aichengWebpage.SavePath)
		if !exists{
			os.MkdirAll(aichengWebpage.SavePath,os.ModePerm)
		}
		go aichengWebpage.ParserWebpage(targetUrlStrSub,aichengHost,client)

		for i:= 0;i<aichengWebpage.CurrentTask;i++{
			go 	aichengWebpage.DownloadPicAndSeed(client,out)
		}
		temp := 0
		fmt.Println("Download has started ,please waitting...")
forEnd:		for  {
/*				temp += <-out
				if 	(aichengWebpage.ItemNum == temp){
					fmt.Printf("共收获%d个条目\n",aichengWebpage.ItemNum)
					break forEnd
				}*/
			select {
				case <-out:
					temp++
					if temp == aichengWebpage.ItemNum{
						fmt.Printf("共收获%d个条目\n",aichengWebpage.ItemNum)

						break forEnd
					}
				case <-time.After(300*time.Second):
						fmt.Println("程序终结于超时(time out)\n")
					break forEnd
			}
		}

	}else{
		caoliuWebpage := parserWebpage.CaoliuWebpage{cmdLineOption,seedAndPic,0}
		caoliuWebpage.SavePath = downloadDirName
		//		fmt.Println(aichengWebpage.SavePath)
		exists,_ := utils.PathExists(caoliuWebpage.SavePath)
		if !exists{
			os.MkdirAll(caoliuWebpage.SavePath,os.ModePerm)
		}
		go caoliuWebpage.ParserWebpage(targetUrlStrSub,caoliuHost,client)

		for i:= 0;i<caoliuWebpage.CurrentTask;i++{
			go 	caoliuWebpage.DownloadPicAndSeed(client,out)
		}
		temp := 0
		fmt.Println("Download has started ,please waitting...")
	forEndd:		for  {
		/*				temp += <-out
						if 	(caoliuWebpage.ItemNum == temp){
							fmt.Printf("共收获%d个条目\n",caoliuWebpage.ItemNum)
							break forEnd
						}*/
		select {
		case <-out:
			temp++
			if temp == caoliuWebpage.ItemNum{
				fmt.Printf("共收获%d个条目\n",caoliuWebpage.ItemNum)
				break forEndd
			}
		case <-time.After(300*time.Second):
			fmt.Println("程序终结于超时(time out)\n")
			break forEndd
		}
	}
	}


}

func showCmdOrder(order cmdInteracion.CmdOrder){
	fmt.Printf("Save-path :%s \nCurrent-task :%d \nAv-class :%s \n",order.SavePath,order.CurrentTask,order.AvClass)
	fmt.Printf("Item range from %d to %d \nProxy :%s\n",order.TopicRange[0],order.TopicRange[1],order.ProxyString)
	fmt.Printf("The love key words:")
	for v:=range order.TopicLike  {
		fmt.Printf("  %s",v)
	}
	fmt.Println()
	fmt.Printf("The hated key words:")
	for v:=range order.TopicHate  {
		fmt.Printf("  %s",v)
	}
	fmt.Println()
}
//todo solve the log error : Unsolicited response received on idle HTTP channel starting with "\n"; err=<nil>
//想要解决这个问题估计得要了解一下《http详解》和GO的开发包
func InitClient(proxyStr string) *http.Client{
	var client *http.Client
	urli := url.URL{}
	urlProxy,_:= urli.Parse(proxyStr)
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
				ResponseHeaderTimeout:30*time.Second,
			},
//			Timeout:30*time.Second,
		}
	}
	return client
}


//todo 为何如下的函数不可以初始化client
/*func InitClient(client *http.Client,proxyStr string) {
	urli := url.URL{}
	urlProxy,_:= urli.Parse(proxyStr)
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	}
}*/


func getParserTargetUrl(avClass string) (string,bool,error) {
	var webpageSite string
	var aichengFlag bool
	err := errors.New("")
	err = nil
	switch avClass {
	case "caoliu_west_original":
		webpageSite = "thread0806.php?fid=4"
		aichengFlag =false
	case "caoliu_cartoon_original":
		webpageSite ="thread0806.php?fid=5"
		aichengFlag = false
	case "caoliu_asia_mosaicked_original":
		webpageSite = "thread0806.php?fid=15"
		aichengFlag = false
	case "caoliu_asia_non_mosaicked_original":
		webpageSite = "thread0806.php?fid=2"
		aichengFlag = false
/*	case "caoliu_west_reposted":
		webpageSite ="thread0806.php?fid=19"
		aichengFlag = false*/
/*	case "caoliu_cartoon_reposted":
		webpageSite = "thread0806.php?fid=24"
		aichengFlag = false*/
/*	case "caoliu_asia_non_mosaicked_reposted":
		webpageSite = "thread0806.php?fid=17"
		aichengFlag = false*/
/*	case "caoliu_asia_mosaicked_reposted":
		webpageSite ="thread0806.php?fid=1"
		aichengFlag = false*/
	case "caoliu_selfie":
		webpageSite ="thread0806.php?fid=16"
		aichengFlag = false
	case "aicheng_west":
		webpageSite= "thread.php?fid=5"
		aichengFlag = true
	case "aicheng_cartoon":
		webpageSite = "thread.php?fid=6"
		aichengFlag = true
	case "aicheng_asia_mosaicked":
		webpageSite="thread.php?fid=4"
		aichengFlag = true
	case "aicheng_asia_non_mosaicked":
		webpageSite="thread.php?fid=16"
		aichengFlag= true
	case "aicheng_selfie":
		webpageSite = "thread.php?fid=22"
		aichengFlag = true
	default:
		err	=errors.New("Wrong av-class!")
	}
	return webpageSite,aichengFlag,err
}