package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
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

	var msgId uint16 = 101

	m := md5.New()
	m.Write([]byte("123456"))
	var request map[string]string
	request = make(map[string]string)
	request["account"] = "ella"
	request["pwd"] = hex.EncodeToString(m.Sum(nil))

	reqBytes,_ := json.Marshal(request)

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
	json.Unmarshal(body,&obj)
	objparams := obj.(map[string]interface{})
	var params map[string]string
	params = make(map[string]string)
	for k,v := range objparams {
		params[k] = v.(string)
	}
    fmt.Println(params)

    msgId = 0
    var request2 map[string]string
    request2 = make(map[string]string)
	request2["verify"] = params["token"]
	reqBytes2,_ := json.Marshal(request2)

	var reqBuff2 []byte = make([]byte, 4+len(reqBytes2))

	binary.BigEndian.PutUint16(reqBuff2[0:2], uint16(len(reqBuff2)))
	binary.BigEndian.PutUint16(reqBuff2[2:4], msgId)
	copy(reqBuff2[4:], reqBytes2)

	// ping <--> pong
	// write
	conn.Write(reqBuff2)
	// read
	p2, err := protocol.ReadPacket(conn)
	checkError(err)

	body2 := p2.GetBody()
	var obj2 interface{}
	json.Unmarshal(body2,&obj2)
	objparams2 := obj2.(map[string]interface{})
	var params2 map[string]string
	params2 = make(map[string]string)
	for k,v := range objparams2 {
		params2[k] = v.(string)
	}
    fmt.Println(params2)

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
