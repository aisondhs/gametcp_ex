import socket
import struct
import json

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('127.0.0.1', 8989))


msgId = 100;  #Signup
request ={}
request["account"] = "ella"
request["pwd"] = "123456"
request["srvid"] = "1"
reqContent = json.dumps(request)

data = struct.pack(">H",len(reqContent)+4)+struct.pack(">H",msgId)+reqContent
sock.send(data)
res = sock.recv(1024)
#print struct.unpack(">H",res[0:2])  #lenghth
#print struct.unpack(">H",res[2:4])  #msgId
rspdata = json.loads(res[4:])
print rspdata


msgId = 101;  #Login
#request ={}
#request["account"] = "ella"
#request["pwd"] = "123456"
#request["srvid"] = "1"

data = struct.pack(">H",len(reqContent)+4)+struct.pack(">H",msgId)+reqContent
sock.send(data)
res = sock.recv(1024)
#print struct.unpack(">H",res[0:2])  #lenghth
#print struct.unpack(">H",res[2:4])  #msgId
rspdata = json.loads(res[4:])
print rspdata

msgId = 0;  # test verify
request ={}
request["verify"] = rspdata["token"]
reqContent = json.dumps(request)

data = struct.pack(">H",len(reqContent)+4)+struct.pack(">H",msgId)+reqContent
sock.send(data)
res = sock.recv(1024)
rspdata = json.loads(res[4:])
print rspdata

sock.close()