package call

import (
	"log"
	"time"

	"github.com/tarm/serial"
	"github.com/vadimipatov/gcircularqueue"
)

func Init(serial_name string, serial_baud int) {
	var err error

	//配置串口读取队列
	Queue = gcircularqueue.NewCircularQueue(10240)
	c := &serial.Config{Name: serial_name, Baud: serial_baud}
	Modem_Port, err = serial.OpenPort(c)

	if err != nil {
		log.Fatal(serial_name + ":打开失败." + err.Error())
	}
	Modem_Port.Flush()
	go Reading_byte()
	time.Sleep(1 * time.Second)

	Modem_Port.Flush()
}
