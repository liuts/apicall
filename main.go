package main

import (
	"apicall/call"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"net/http"
	"os"
)

//电话队列长度设置为2
var Phone_chan chan string = make(chan string, 2)

func main() {
	//读取配置文件
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// 初始化串口
	serial_name := cfg.Section("").Key("serial").String() //  /dev/ttyS1
	serial_baud, _ := cfg.Section("").Key("baud").Int()
	call.Init(serial_name, serial_baud)
	defer call.Modem_Port.Close()

	//开启打电话go协程
	go go_call()

	//初始化http服务
	http_port := cfg.Section("").Key("http").String()
	router := gin.Default()
	router.GET("/call", handler_call)
	router.Run(":" + http_port)
}

func go_call() {
	for v := range Phone_chan {
		call.Make_call(v)
	}
}

//匹配的url格式:  /call?phone=138138138138&key=iamops
func handler_call(c *gin.Context) {
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
		Phone_chan <- phone
	}
}
