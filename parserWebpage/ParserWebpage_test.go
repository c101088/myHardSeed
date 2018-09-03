package parserWebpage

import (
	"testing"
	"net/http"
	"net/url"
	"fmt"
	"self-goProject/myHardSeed/cmdInteracion"
	"github.com/godfather1103/hardseed-go/utils"
	"os"
)

func TestForTestingRegx(t *testing.T) {

		ForTestingRegx()
}

func TestAichengWebpage_ParserWebpage(t *testing.T) {
/*	ichengWebpage := AichengWebpage{cmdInteracion.CmdOrder{
		"C:/",
		16,
		"aicheng_asia_mosaicked",
		[]int{0,40},
		"socks5://127.0.0.1:1080",
		nil,
		nil,
		},
		make(chan SeedAndPicChan),
	}
*/
	var client *http.Client
	urli := url.URL{}
	urlProxy,_:= urli.Parse("socks5://127.0.0.1:1080")
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	}
	_, err := client.Get("http://www.ac168.info/bt/thread.php?fid=16")
	if err != nil  {
		fmt.Printf("aaaa\n")
	}
	//ichengWebpage.ParserWebpage("thread.php?fid=16","http://www.ac168.info/bt/",client)
}

func TestAichengWebpage_DownloadPicAndSeed(t *testing.T) {
	ichengWebpage := AichengWebpage{cmdInteracion.CmdOrder{
		"C:\\download\\",
		16,
		"aicheng_asia_mosaicked",
		[]int{0,40},
		"socks5://127.0.0.1:1080",
		nil,
		nil,
	},
		make(chan SeedAndPicChan),
	10,
	}

	exists,_ := utils.PathExists(ichengWebpage.SavePath)
	if !exists{
		os.Mkdir(ichengWebpage.SavePath,os.ModePerm)
	}

	picUrl := []string{"http://imagizer.imageshack.com/img923/6171/1jgRpm.jpg","http://imagizer.imageshack.com/img922/3054/1wYHy7.jpg"}


	var client *http.Client
	urli := url.URL{}
	urlProxy,_:= urli.Parse("socks5://127.0.0.1:1080")
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	}

}
