package main

import (
	"apicall/call"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"net/http"
	"os"
)

// PhoneChan 电话队列长度设置为2000
var PhoneChan chan string = make(chan string, 2000)

func main() {
	//读取配置文件
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// 初始化串口
	serialName := cfg.Section("").Key("serial").String() //  /dev/ttyS1
	serialBaud, _ := cfg.Section("").Key("baud").Int()
	call.Init(serialName, serialBaud)
	defer call.Modem_Port.Close()

	//开启打电话go协程
	go goCall()

	//初始化http服务
	httpPort := cfg.Section("").Key("http").String()
	router := gin.Default()
	router.GET("/call", handlerCall)
	router.GET("/clean", handlerClean)

	router.Run(":" + httpPort)
}

func goCall() {
	for v := range PhoneChan {
		call.Make_call(v)
		//fmt.Printf("call=%s\n", v)
		//time.Sleep(2 * time.Second)
	}
}

func handlerClean(c *gin.Context) {
	var phones string
	var countPhones string
	count := len(PhoneChan)
	for len(PhoneChan) > 0 {
		phones = <-PhoneChan
		countPhones = countPhones + phones + ","
	}
	c.JSON(http.StatusOK, gin.H{"message": "clean=" + fmt.Sprintf("%d", count) + "|" + countPhones})
}

// 匹配的url格式:  /call?phone=138138138138&key=iamops
func handlerCall(c *gin.Context) {
	phone := c.Query("phone")
	key := c.Query("key") // 是 c.Request.URL.Query().Get("lastname") 的简写

	if key != "iamops" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Failed",
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK," + phone + key,
		})
		PhoneChan <- phone
	}
}
