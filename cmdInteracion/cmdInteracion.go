package cmdInteracion

import (
	"flag"
	"runtime"
	"strings"
	"strconv"
)

type CmdOrder struct {
	SavePath string
	CurrentTask int
	AvClass string
	TopicRange []int
	ProxyString string
	TopicLike []string
	TopicHate []string
	CoreNum int
}

func ParserCmd(cmdInstance *CmdOrder)( error){
	var strTopicRange string
	var strTopicLike string
	var strTopicHate string
	flag.StringVar(&cmdInstance.SavePath,"save-path",sysClass(),"下载文件的保存路径 ")
	flag.IntVar(&cmdInstance.CurrentTask,"current-task",16,"下载的并行任务数")
	flag.StringVar(&cmdInstance.AvClass,"av-class","caoliu_asia_non_mosaicked_original","种子分区")
	flag.StringVar(&strTopicRange,"topic-range","0 10","文件范围")
	flag.StringVar(&cmdInstance.ProxyString,"proxy","socks5://127.0.0.1:1080","代理端口，默认为ss")
	flag.StringVar(&strTopicLike,"like","","钟爱的主题，如ABP？")
	flag.IntVar(&cmdInstance.CoreNum,"core-number",runtime.NumCPU(),"The cpu core used for hardseed , more than 2 core would be welcome!")

	//TODO 剔除“连发 合集”等主题
	flag.StringVar(&strTopicHate,"hate","","憎恨的主题，如重口？")
	flag.Parse()

	for _,v :=range strings.Split(strTopicRange," "){
		intValue,_ := strconv.Atoi(v)
		cmdInstance.TopicRange= append(cmdInstance.TopicRange,intValue)
	}
	//空字符穿分割也会形成有内容的[]string，故处理
	if strTopicLike ==""{
		cmdInstance.TopicLike = nil
	}else {
		for _,v :=range strings.Split(strTopicLike," "){
			cmdInstance.TopicLike= append(cmdInstance.TopicLike,v)
		}

	}

	if strTopicHate ==""{
		cmdInstance.TopicHate = nil
	}else{
		for _,v :=range strings.Split(strTopicHate," "){
			cmdInstance.TopicHate= append(cmdInstance.TopicHate,v)
		}
	}
	return nil
}

func sysClass() string{
	switch runtime.GOOS {
	case "darwin" :
		return string("/user/")
	case "linux":
		return string( "~/")
	case "windows":
		return string("D:\\download\\")
	default:
		return "unknow opeartion system!"
	}
}