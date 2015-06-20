package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/aisondhs/gametcp_ex/protocol"
	//"io/ioutil"
	"log"
	"net"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	var msgId uint16 = 100 // Sign up
	var request map[string]string
	request = make(map[string]string)
	request["account"] = "ella"
	request["pwd"] = "123456"
	request["srvid"] = "1"

	reqBytes, _ := json.Marshal(request)

	var reqBuff []byte = make([]byte, 4+len(reqBytes))

	binary.BigEndian.PutUint16(reqBuff[0:2], uint16(len(reqBuff)))
	binary.BigEndian.PutUint16(reqBuff[2:4], msgId)
	copy(reqBuff[4:], reqBytes)

	// ping <--> pong
	// write
	conn.Write(reqBuff)
	// read
	p, err := protocol.ReadPacket(conn)
	checkError(err)

	body := p.GetBody()
	var obj interface{}
	json.Unmarshal(body, &obj)
	objparams := obj.(map[string]interface{})
	var params map[string]string
	params = make(map[string]string)
	for k, v := range objparams {
		params[k] = v.(string)
	}
	fmt.Println(params)

	msgId = 101 // login

	//var request map[string]string
	//request = make(map[string]string)
	//request["account"] = "ella"
	//request["pwd"] = "123456"
	//request["Areaid"] = 1

	//reqBytes, _ := json.Marshal(request)

	var reqBuff2 []byte = make([]byte, 4+len(reqBytes))

	binary.BigEndian.PutUint16(reqBuff2[0:2], uint16(len(reqBuff2)))
	binary.BigEndian.PutUint16(reqBuff2[2:4], msgId)
	copy(reqBuff2[4:], reqBytes)

	// ping <--> pong
	// write
	conn.Write(reqBuff2)
	// read
	p2, err := protocol.ReadPacket(conn)
	checkError(err)

	body2 := p2.GetBody()
	var obj2 interface{}
	json.Unmarshal(body2, &obj2)
	objparams2 := obj2.(map[string]interface{})
	var params2 map[string]string
	params2 = make(map[string]string)
	for k, v := range objparams2 {
		params2[k] = v.(string)
	}
	fmt.Println(params2)

	msgId = 0 // test verify
	var request3 map[string]string
	request3 = make(map[string]string)
	request3["verify"] = params2["token"]
	reqBytes3, _ := json.Marshal(request3)

	var reqBuff3 []byte = make([]byte, 4+len(reqBytes3))

	binary.BigEndian.PutUint16(reqBuff3[0:2], uint16(len(reqBuff3)))
	binary.BigEndian.PutUint16(reqBuff3[2:4], msgId)
	copy(reqBuff3[4:], reqBytes3)

	// ping <--> pong
	// write
	conn.Write(reqBuff3)
	// read
	p3, err := protocol.ReadPacket(conn)
	checkError(err)

	body3 := p3.GetBody()
	var obj3 interface{}
	json.Unmarshal(body3, &obj3)
	objparams3 := obj3.(map[string]interface{})
	var params3 map[string]string
	params3 = make(map[string]string)
	for k, v := range objparams3 {
		params3[k] = v.(string)
	}
	fmt.Println(params3)

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
