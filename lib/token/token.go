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
	Get(key string) (string, error)
	Del(key string) error
	Expire(key string, time int64) (bool, error)
}

type Token struct {
	adapter  adapter
	isUnique bool
}

func NewToken(a adapter) *Token {
	return &Token{a, true}
}

func (this *Token) NotUnique() {
	this.isUnique = false
}

// get uid from token
func (this *Token) GetUid(token string) (int64, error) {

	if str, err := this.adapter.Get(token); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(string(str), 10, 0)
	}
}

// create new token
func (this *Token) AddToken(uid int64) (string, error) {

	m := md5.New()
	m.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(uid, 10)))
	token := hex.EncodeToString(m.Sum(nil))

	if this.isUnique {
		this.setUidToken(uid, token)
	}
	err := this.adapter.Set(token, strconv.Itoa(int(uid)))
	if err != nil {
		return "",err
	}
	this.SetExpire(token,7200)
	return token,err
}

func (this *Token) setUidToken(uid int64, token string) error {

	key := "uid_token_" + strconv.Itoa(int(uid))

	if oldToken, err := this.adapter.Get(key); err == nil {
		this.adapter.Del(string(oldToken))
	}

	return this.adapter.Set(key, token)
}

func (this *Token) SetExpire(key string,time int64) error {
	_,err := this.adapter.Expire(key, time)
	return err
}
