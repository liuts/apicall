# apicall


使用golang实现的使用sim800模块打电话的接口,用来作为zabbix的报警接口.
运行流程:程序启动后,http接收到电话号码后,会放到一个channel当中,打电话的协程会从channel中取出
电话号码,调用串口子程序打电话.串口的收发使用github.com/vadimipatov/gcircularqueue做环形队列.减少因干扰造成数据
丢失,影响程序运行的可能.

#### 使用方法:

### 一 需要准备的模块:

![image](https://raw.githubusercontent.com/liuts/apicall/master/sim800.png?raw=true)
##### 1,准备一个sim800芯片的模组,我用的是"USB转GSM串口GPRSSIM800C模块"

##### 2,一张GSM的电话卡

##### 3,物理机一台,安装centos7 或者 windows10

### 二 调试

##### 1,将GSM电话卡插入模块中,模块开机,试用命令  ls -l /dev/tty*,将正确的串口号配置在程序的config.ini文件中

##### 2,直接go build 启动.运行生成的文件,同时观察日志输出.

##### 3,程序运行后会运行在8181端口,之后直接在浏览器中访问http://127.0.0.1:8181/call?phone=138138138138&key=iamops 

#### todo：
- ✅1,增加一个状态接口,用来查看当前队列长度和队列内容
- 2,增加一个配置文件,用来配置串口号,波特率,超时时间等
- 3,增加一个日志文件,用来记录程序运行日志
