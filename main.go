package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/aisondhs/gametcp_ex/controllers"
	"github.com/aisondhs/gametcp_ex/lib/funcmap"
	"github.com/aisondhs/gametcp_ex/lib/gametcp"
	"github.com/aisondhs/alog"
	"github.com/aisondhs/gametcp_ex/protocol"
)

type Callback struct{}

var funcs funcmap.Funcs

var HTTP_PORT string

var logdir string = "./logs"

var actList map[uint16](string)

func init() {
	// bind func map
	funcs = funcmap.NewFuncs(100)
	actList = make(map[uint16](string), 100)

	funcs.Bind("Hello", controllers.Hello)
	funcs.Bind("Login", controllers.Login)
	funcs.Bind("Signup", controllers.Signup)
	actList[0] = "Hello"
	actList[100] = "Signup"
	actList[101] = "Login"
	alog.Init(logdir,alog.ROTATE_BY_DAY,false)
}

func (this *Callback) OnConnect(c *gametcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	alog.Info("OnConnect:"+addr.String())
	return true
}

func (this *Callback) OnMessage(c *gametcp.Conn, p protocol.Packet) bool {
	packet := &p

	reqContent := packet.GetBody()
	msgId := packet.GetMsgId()
	var obj interface{}
	json.Unmarshal(reqContent, &obj)
	objparams := obj.(map[string]interface{})
	var params map[string]string
	params = make(map[string]string)
	for k, v := range objparams {
		params[k] = v.(string)
	}
	methodName := actList[msgId]

	var response map[string]string
	response = make(map[string]string)

	var rid string

	if methodName != "Login" && methodName != "Signup" {
		//verify
		verifyInfo, err := controllers.Verify(params["verify"])
		if err != nil {
			response["status"] = "fail"
			response["msg"] = err.Error()
		} else {
			reflectData, err := funcs.Call(methodName, params)
			checkError(err)
			i := reflectData[0].Interface()
			response = i.(map[string]string)
		}
		rid = verifyInfo["rid"]
	} else {
		reflectData, err := funcs.Call(methodName, params)
		checkError(err)
		i := reflectData[0].Interface()
		response = i.(map[string]string)
		rid = response["rid"]
	}
	rspBytes, _ := json.Marshal(response)
	rspPacket := protocol.NewPacket(rspBytes, msgId, false)
	c.AsyncWritePacket(rspPacket, time.Second)

	logData := alog.Mrecord{"rid":rid,"req":params,"rsp":response}
	alog.Info(methodName,logData)

	return true
}

func (this *Callback) OnClose(c *gametcp.Conn) {
	alog.Info("OnClose:"+c.GetExtraData().(net.Addr).String())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		log.Fatal(err)
	}

	HTTP_PORT, err = c.GetValue("Server", "port")
	checkError(err)

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+HTTP_PORT)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	sendChan, err := c.Int("Server", "sendChan")
	checkError(err)
	receiveChan, err := c.Int("Server", "receiveChan")
	checkError(err)

	// creates a server
	config := &gametcp.Config{
		PacketSendChanLimit:    uint32(sendChan),
		PacketReceiveChanLimit: uint32(receiveChan),
	}
	srv := gametcp.NewServer(config, &Callback{})

	// starts service
	go srv.Start(listener, time.Second*5)
	alog.Info("listening:"+listener.Addr().String())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-chSig
	alog.Error("listening:"+sig.String())

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		alog.Error(err.Error())
	}
}
