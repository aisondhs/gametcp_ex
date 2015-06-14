package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"encoding/json"
	//"fmt"

	"github.com/Unknwon/goconfig"
	"github.com/aisondhs/gametcp_ex/controllers"
	"github.com/aisondhs/gametcp_ex/lib/funcmap"
	"github.com/aisondhs/gametcp_ex/lib/gametcp"
	"github.com/aisondhs/gametcp_ex/lib/logger"
	"github.com/aisondhs/gametcp_ex/protocol"
)

type Callback struct{}

var funcs funcmap.Funcs

var HTTP_PORT string

var logdir string = "./logs"

var actList map[uint16](string)

func init() {
	// 返回数据。
	funcs = funcmap.NewFuncs(100)
	actList = make(map[uint16](string), 100)

	funcs.Bind("Hello", controllers.Hello)
	funcs.Bind("Login", controllers.Login)
	funcs.Bind("Signup", controllers.Signup)
	actList[0] = "Hello"
	actList[100] = "Signup"
	actList[101] = "Login"
}

func (this *Callback) OnConnect(c *gametcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	logger.PutLog("OnConnect:"+addr.String(), logdir, "info")
	return true
}

func (this *Callback) OnMessage(c *gametcp.Conn, p protocol.Packet) bool {
	packet := &p

	reqContent := packet.GetBody()
	msgId := packet.GetMsgId()
	var obj interface{}
	json.Unmarshal(reqContent,&obj)
	objparams := obj.(map[string]interface{})
	var params map[string]string
	params = make(map[string]string)
	for k,v := range objparams {
		params[k] = v.(string)
	}
	methodName := actList[msgId]

	var response map[string]string
	response = make(map[string]string)

	if methodName != "Login" && methodName != "Signup" {
		//verify
		_,err := controllers.Verify(params["verify"])
		if err != nil {
			response["status"] = "fail"
			response["msg"] = "Please Login first! "+err.Error()
		} else {
			reflectData, err := funcs.Call(methodName, params)
	        if err != nil {
		        logger.PutLog(err.Error(), logdir, "error")
	        }
		    i := reflectData[0].Interface()
		    response = i.(map[string]string)
		}
	} else {
		reflectData, err := funcs.Call(methodName, params)
	    if err != nil {
		    logger.PutLog(err.Error(), logdir, "error")
	    }
		i := reflectData[0].Interface()
		response = i.(map[string]string)
	}
	rspBytes,_ := json.Marshal(response)
	rspPacket := protocol.NewPacket(rspBytes, msgId, false)
	c.AsyncWritePacket(rspPacket, time.Second)
	logmsg := "Req: " + string(reqContent) + " Rsp: " + string(rspBytes)
	logger.PutLog(logmsg, logdir, methodName)
	return true
}

func (this *Callback) OnClose(c *gametcp.Conn) {
	logger.PutLog("OnClose:"+c.GetExtraData().(net.Addr).String(), logdir, "info")
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
	go srv.Start(listener, time.Second)
	logger.PutLog("listening:"+listener.Addr().String(), logdir, "info")

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-chSig
	logger.PutLog("listening:"+sig.String(), logdir, "error")

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		logger.PutLog(err.Error(), logdir, "error")
	}
}
