package token

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/aisondhs/gametcp_ex/lib/redis"
	"strconv"
	"time"
)

var GameToken = NewToken(redis.Redis)

func init() {
}

type adapter interface {
	Set(key string, value string) error
	Hmset(key string, params map[string](string)) error
	Hmget(key string, params []string) (map[string](string), error)
	Hgetall(key string) (map[string](string), error)
	Get(key string) (string, error)
	Del(key string) error
	Expire(key string, time int64) (bool, error)
}

type Token struct {
	adapter adapter
}

func NewToken(a adapter) *Token {
	return &Token{a}
}

// get uid from token
func (this *Token) GetTokenInfo(token string) (map[string]string, error) {
	keys := []string{"uid", "rid", "srvid"}
	return this.adapter.Hmget(token, keys)
}

// create new token
func (this *Token) AddToken(uid int64, rid int64, srvid int32) (map[string]string, error) {
	m := md5.New()
	m.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(uid, 10) + strconv.FormatInt(rid, 10)))
	token := hex.EncodeToString(m.Sum(nil))

	var params map[string]string
	params = make(map[string]string)
	params["uid"] = strconv.Itoa(int(uid))
	params["rid"] = strconv.Itoa(int(rid))
	params["srvid"] = strconv.Itoa(int(srvid))

	err := this.adapter.Hmset(token, params)
	if err != nil {
		return nil, err
	}
	//delete(params,"areaId")
	params["lastlogin"] = strconv.Itoa(int(time.Now().Unix()))
	params["token"] = token
	this.SetUidInfo(uid, rid, params)
	this.SetExpire(token, 7200)
	return params, nil
}

func (this *Token) SetUidInfo(uid int64, rid int64, params map[string]string) error {
	key := "token_" + strconv.Itoa(int(uid)) + "_" + strconv.Itoa(int(rid))
	uidInfo, _ := this.GetUidInfo(uid, rid)
	if oldToken, ok := uidInfo["token"]; ok {
		this.adapter.Del(oldToken)
	}
	return this.adapter.Hmset(key, params)
}

func (this *Token) GetUidInfo(uid int64, rid int64) (map[string]string, error) {
	key := "token_" + strconv.Itoa(int(uid)) + "_" + strconv.Itoa(int(rid))
	return this.adapter.Hgetall(key)
}

func (this *Token) SetExpire(key string, time int64) error {
	_, err := this.adapter.Expire(key, time)
	return err
}
