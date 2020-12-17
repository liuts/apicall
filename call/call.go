package call

import (
	"github.com/tarm/serial"
	"github.com/vadimipatov/gcircularqueue"
	"log"
	"strings"
	"time"
)

var Modem_Port *serial.Port
var Queue *gcircularqueue.CircularQueue

func Reading_byte() {
	for true {
		time.Sleep(1 * time.Millisecond)
		buf := make([]byte, 1)
		_, err := Modem_Port.Read(buf)
		if err != nil {
			log.Println(err)

		} else {
			Queue.Push(buf[0])
			//fmt.Println(buf[0])
		}
	}

}

func Get_reply_word() string {
	//2020/12/02 17:01:01 "AT\r\nOK\r\n"
	n := Queue.Len()
	words := ""

	for i := 0; i < n; i++ {
		tmp := Queue.Shift()
		if tmp != nil {
			// \r=13,\n=10
			if tmp.(byte) == 10 {
				log.Println("<--", strings.TrimSpace(words))
				return words
			} else {
				words = words + string(tmp.(byte))
			}
		} else {
			return words
		}

	}
	return words

}

//sendCommand(char *Command, char *Response, unsigned long Timeout, unsigned char Retry)
func sendCommand(Command string, Response string, Timeout int, Retry int) bool {
	Timeout = Timeout * 10
	for i := 0; i < Retry; i++ {
		Modem_WriteLine(Command)

		for t := 0; t < Timeout; t++ {
			time.Sleep(100 * time.Millisecond) //0.01second
			bb := Get_reply_word()
			//fmt.Println(strings.TrimSpace(bb) , strings.TrimSpace(Response))
			if strings.TrimSpace(bb) == strings.TrimSpace(Response) {
				return true
			}

		}
	}
	return false
}

func Modem_WriteLine(line string) {
	log.Println("-->", line)
	_, err := Modem_Port.Write([]byte(line + "\n"))
	if err != nil {
		log.Println(err)
	}
}

func Make_call(phone string) {
	ok := sendCommand("AT", "OK", 10, 2)
	//Modem_Port.Flush()

	ok = sendCommand("ATH", "OK", 10, 2)
	//Modem_Port.Flush()

	//AT+MORING=1

	ok = sendCommand("AT+MORING=1", "OK", 10, 2)
	//Modem_Port.Flush()

	ok = sendCommand("ATD"+phone+";", "MO CONNECTED", 100, 1)
	//Modem_Port.Flush()
	if ok != true {
		ok = sendCommand("ATH", "OK", 10, 2)
		//Modem_Port.Flush()
		log.Println("zhi xing jie guo", ok, Queue.Len())
	}
	time.Sleep(1 * time.Second)
	ok = sendCommand("AT+CTTS=2,\"import message\"", "+CTTS: 0", 100, 2)
	//Modem_Port.Flush()

	time.Sleep(1 * time.Second)
	ok = sendCommand("AT+CTTS=2,\"import message\"", "+CTTS: 0", 100, 2)
	//Modem_Port.Flush()

	time.Sleep(2 * time.Second)
	ok = sendCommand("ATH", "OK", 10, 2)
	Modem_Port.Flush()
	//defer Modem_Port.Close()
}

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
