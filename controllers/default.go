package controllers

import (
	//"crypto/md5"
	// "encoding/hex"
	//"fmt"
	"github.com/aisondhs/gametcp_ex/lib/logger"
	//"github.com/aisondhs/gametcp_ex/lib/redis"
	"github.com/aisondhs/gametcp_ex/lib/token"
	"github.com/aisondhs/gametcp_ex/models"
	//redigo "github.com/garyburd/redigo/redis"
	"strconv"
	//"strings"
	"time"
)

func Verify(verifyKey string) (int64,error){
	uid,err := token.GameToken.GetUid(verifyKey)
	if err != nil {
		return 0,err
	}
	token.GameToken.SetExpire(verifyKey,7200)
	return uid,err
}

func Hello(request map[string](string)) map[string](string) {
	var response map[string]string
	response = make(map[string]string)
	response["status"] = "1"
    response["msg"] = "success"
    return response
}

func Signup(request map[string](string)) map[string](string) {
	account := request["account"]
	pwd := request["pwd"]
	user, err := models.User.GetUserByAccount(account)
	checkError(err)
	var response map[string]string
	response = make(map[string]string)
	if user.Uid == 0 {
		user = new(models.UserTable)
		user.Account = account
		user.Pwd = pwd
		now := time.Now()
		user.Ctime = now.Unix()
		err = models.User.Insert(user)
		checkError(err)

		//role := new(models.RoleTable)
		role , err := models.Role.Insert(user.Uid,1)
		if err != nil {
			response["status"] = "0"
			response["msg"] = err.Error()
		} else {
			response["status"] = "1"
		    response["msg"] = "Signup success,uid :" + strconv.Itoa(int(user.Uid))+" role id:"+strconv.Itoa(int(role.RoleId))
		}
		
	} else {
		response["status"] = "0"
		response["msg"] = "Signup fail,account is exists"
	}
	return response
}

func Login(request map[string](string)) map[string](string) {
	account := request["account"]
	pwd := request["pwd"]

	user, err := models.User.GetUserByAccount(account)
	checkError(err)
	var response map[string]string
	response = make(map[string]string)
	if user.Uid == 0 {
		response["status"] = "0"
		response["msg"] = "use not exists"
		response["token"] = ""
	} else {
		if user.Pwd == pwd {
			tokenStr, err := token.GameToken.AddToken(user.Uid)
			checkError(err)
			response["status"] = "1"
		    response["msg"] = "success"
		    response["token"] = tokenStr
		} else {
			response["status"] = "0"
		    response["msg"] = "account or passwd is incollect"
		    response["token"] = ""
		}
	}
	return response
}

func checkError(err error) {
	if err != nil {
		logger.PutLog(err.Error(), "./logs", "error")
	}
}
