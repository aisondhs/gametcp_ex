package controllers

import (
	"crypto/md5"
	"encoding/hex"
	//"fmt"
	//"github.com/aisondhs/gametcp_ex/lib/redis"
	"github.com/aisondhs/gametcp_ex/lib/token"
	"github.com/aisondhs/gametcp_ex/models"
	//redigo "github.com/garyburd/redigo/redis"
	"strconv"
	//"strings"
	"errors"
	"time"
)

func Verify(verifyKey string) (map[string]string, error) {
	tokenInfo, err := token.GameToken.GetTokenInfo(verifyKey)
	if err != nil {
		return nil, err
	}
	if _, ok := tokenInfo["uid"]; ok == false {
		return nil, errors.New("token not found")
	}
	uid, _ := strconv.Atoi(tokenInfo["uid"])
	rid, _ := strconv.Atoi(tokenInfo["rid"])
	userInfo, err := token.GameToken.GetUidInfo(int64(uid), int64(rid))
	if err != nil {
		return nil, errors.New("token not found")
	}
	if verifyKey != userInfo["token"] {
		token.GameToken.SetExpire(verifyKey, 0)
		return nil, errors.New("user multi login")
	}
	token.GameToken.SetExpire(verifyKey, 7200)
	return tokenInfo, nil
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
	m := md5.New()
	m.Write([]byte(pwd))
	password := hex.EncodeToString(m.Sum(nil))
	srvid := request["srvid"]

	var response map[string]string
	response = make(map[string]string)

	user, err := models.User.GetUserByAccount(account)
	if err != nil {
		response["status"] = "0"
		response["msg"] = err.Error()
	} else {
		if user.Uid == 0 {
			user = new(models.UserTable)
			user.Account = account
			user.Pwd = password
			now := time.Now()
			user.Ctime = now.Unix()
			err = models.User.Insert(user)
			if err != nil {
				response["status"] = "0"
				response["msg"] = err.Error()
			} else {
				sid, _ := strconv.Atoi(srvid)
				role, err := models.Role.Insert(user.Uid, int32(sid))
				if err != nil {
					response["status"] = "0"
					response["msg"] = err.Error()
				} else {
					response["status"] = "1"
					response["msg"] = "Signup success,uid :" + strconv.Itoa(int(user.Uid)) + " role id:" + strconv.Itoa(int(role.RoleId))

					var params map[string]string
					params = make(map[string]string)
					params["uid"] = strconv.Itoa(int(user.Uid))
					params["rid"] = strconv.Itoa(int(role.RoleId))
					params["pwd"] = password
					params["srvid"] = srvid
					params["ctime"] = strconv.Itoa(int(time.Now().Unix()))
					token.GameToken.SetUidInfo(user.Uid, role.RoleId, params)
				}
			}
		} else {
			response["status"] = "0"
			response["msg"] = "Signup fail,account is exists"
		}
	}
	return response
}

func Login(request map[string](string)) map[string](string) {
	account := request["account"]
	pwd := request["pwd"]
	m := md5.New()
	m.Write([]byte(pwd))
	password := hex.EncodeToString(m.Sum(nil))
	srvidStr := request["srvid"]

	var response map[string]string
	response = make(map[string]string)

	user, err := models.User.GetUserByAccount(account)
	if err != nil {
		response["status"] = "0"
		response["msg"] = err.Error()
		response["token"] = ""
	} else {
		if user.Uid == 0 {
			response["status"] = "0"
			response["msg"] = "use not exists"
			response["token"] = ""
		} else {
			if user.Pwd == password {
				srvid, _ := strconv.Atoi(srvidStr)
				sid := int32(srvid)
				role, err := models.Role.GetRoleByArea(user.Uid, sid)
				if err != nil {
					response["status"] = "0"
					response["msg"] = err.Error()
					response["token"] = ""
				} else {
					tokenInfo, err := token.GameToken.AddToken(user.Uid, role.RoleId, sid)
					if err != nil {
						response["status"] = "0"
						response["msg"] = err.Error()
						response["token"] = ""
						response["uid"] = ""
						response["rid"] = ""
					} else {
						response["status"] = "1"
						response["msg"] = "success"
						response["uid"] = strconv.Itoa(int(user.Uid))
						response["rid"] = strconv.Itoa(int(role.RoleId))
						response["token"] = tokenInfo["token"]
					}
				}
			} else {
				response["status"] = "0"
				response["msg"] = "account or passwd is incollect"
				response["token"] = ""
			}
		}
	}
	return response
}
