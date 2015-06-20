package redis

import (
	//"bytes"
	"errors"
	"github.com/Unknwon/goconfig"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

func init() {
	c, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		panic(err)
	}

	ip, err := c.GetValue("Redis", "ip")
	if err != nil {
		panic(err)
	}
	dbnum, err := c.Int("Redis", "dbnum")
	if err != nil {
		panic(err)
	}

	Redis = &myRedis{newPool(ip, dbnum)}
	PoolInitList = make(map[string](*myRedis), 16)
	PoolInitList[ip+strconv.Itoa(dbnum)] = Redis
}

func newPool(server string, dbnum int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", dbnum); err != nil {
				return nil, err
			}

			//			if _, err := c.Do("AUTH", password); err != nil {
			//				c.Close()
			//				return nil, err
			//			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func NewRedisPool(server string, dbnum int) *myRedis {
	mapkey := server + strconv.Itoa(dbnum)
	if newRedis, ok := PoolInitList[mapkey]; ok {
		return newRedis
	}
	newRedis := &myRedis{newPool(server, dbnum)}
	PoolInitList[mapkey] = newRedis
	return newRedis
}

var (
	Redis        *myRedis
	NotFind      = errors.New("Not Find")
	RedisError   = errors.New("Unexpected reply")
	PoolInitList map[string](*myRedis)
)

type myRedis struct {
	*redis.Pool
}

func (r *myRedis) Get(key string) (string, error) {
	conn, err := r.Dial()
	if err != nil {
		return "", err
	}
	result, err := conn.Do("GET", key)
	var s []byte
	if err == nil {
		if result == nil {
			return "", NotFind
		}
		s = result.([]byte)
	}
	return string(s), err
}

func (r *myRedis) Set(key, value string) error {
	conn, err := r.Dial()
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value)
	return err
}

func (r *myRedis) Del(key string) error {
	conn, err := r.Dial()
	if err != nil {
		return err
	}
	_, err = conn.Do("DEL", key)
	return err
}

func (r *myRedis) Type(key string) (string, error) {
	conn, err := r.Dial()
	if err != nil {
		return "", err
	}
	res, err := conn.Do("TYPE", key)

	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func (r *myRedis) Keys(pattern string) ([]string, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(conn.Do("KEYS", pattern))

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *myRedis) Auth(password string) error {
	conn, err := r.Dial()
	if err != nil {
		return err
	}
	_, err = conn.Do("AUTH", password)
	if err != nil {
		return err
	}
	return nil
}

func (r *myRedis) Exists(key string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("EXISTS", key)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Randomkey() (string, error) {
	conn, err := r.Dial()
	if err != nil {
		return "", err
	}
	res, err := redis.Bytes(conn.Do("RANDOMKEY"))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (r *myRedis) Rename(src string, dst string) error {
	conn, err := r.Dial()
	if err != nil {
		return err
	}
	_, err = conn.Do("RENAME", src, dst)
	if err != nil {
		return err
	}
	return nil
}

func (r *myRedis) Renamenx(src string, dst string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("RENAMENX", src, dst)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Dbsize() (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := redis.Int(conn.Do("DBSIZE"))
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (r *myRedis) Expire(key string, time int64) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("EXPIRE", key, strconv.FormatInt(time, 10))

	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Ttl(key string) (int64, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("TTL", key)
	if err != nil {
		return -1, err
	}
	return res.(int64), nil
}

func (r *myRedis) Move(key string, dbnum int) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("MOVE", key, strconv.Itoa(dbnum))
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

func (r *myRedis) Flush(all bool) error {
	conn, err := r.Dial()
	if err != nil {
		return err
	}
	var cmd string
	if all {
		cmd = "FLUSHALL"
	} else {
		cmd = "FLUSHDB"
	}
	_, err = conn.Do(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (r *myRedis) Setnx(key string, val string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}

	res, err := conn.Do("SETNX", key, val)

	if err != nil {
		return false, err
	}
	if data, ok := res.(int64); ok {
		return data == 1, nil
	}
	return false, RedisError
}

func (r *myRedis) Setex(key string, time int64, val string) error {
	conn, err := r.Dial()
	if err != nil {
		return err
	}
	_, err = conn.Do("SETEX", key, strconv.FormatInt(time, 10), val)

	if err != nil {
		return err
	}
	return nil
}

func (r *myRedis) Incr(key string) (int64, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("INCR", key)
	if err != nil {
		return -1, err
	}

	return res.(int64), nil
}

func (r *myRedis) Incrby(key string, val int64) (int64, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("INCRBY", key, strconv.FormatInt(val, 10))
	if err != nil {
		return -1, err
	}
	return res.(int64), nil
}

func (r *myRedis) Decr(key string) (int64, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("DECR", key)
	if err != nil {
		return -1, err
	}
	return res.(int64), nil
}

func (r *myRedis) Decrby(key string, val int64) (int64, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("DECRBY", key, strconv.FormatInt(val, 10))
	if err != nil {
		return -1, err
	}

	return res.(int64), nil
}

func (r *myRedis) Append(key string, val string) (int64, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("APPEND", key, val)

	if err != nil {
		return -1, err
	}
	return res.(int64), err
}

func (r *myRedis) Substr(key string, start int, end int) (string, error) {
	conn, err := r.Dial()
	if err != nil {
		return "", err
	}

	res, err := redis.Bytes(conn.Do("SUBSTR", key, strconv.Itoa(start), strconv.Itoa(end)))
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", RedisError
	}
	return string(res), nil
}

// Set commands

func (r *myRedis) Sadd(key string, value string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("SADD", key, value)

	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Srem(key string, value string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("SREM", key, value)

	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Spop(key string) (string, error) {
	conn, err := r.Dial()
	if err != nil {
		return "", err
	}
	res, err := redis.Bytes(conn.Do("SPOP", key))
	if err != nil {
		return "", err
	}

	if res == nil {
		return "", RedisError
	}
	return string(res), nil
}

func (r *myRedis) Smove(src string, dst string, val string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("SMOVE", src, dst, val)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Scard(key string) (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("SCARD", key)
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func (r *myRedis) Sismember(key string, value string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("SISMEMBER", key, value)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Smembers(key string) ([][]byte, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := conn.Do("SMEMBERS", key)
	if err != nil {
		return nil, err
	}
	return res.([][]byte), nil
}

func (r *myRedis) Srandmember(key string) ([]byte, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := conn.Do("SRANDMEMBER", key)
	if err != nil {
		return nil, err
	}
	return res.([]byte), nil
}

// sorted set commands

func (r *myRedis) Zadd(key string, value string, score float64) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("ZADD", key, strconv.FormatFloat(score, 'f', -1, 64), value)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Zrem(key string, value string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("ZREM", key, value)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Zincrby(key string, value string, score float64) (float64, error) {
	conn, err := r.Dial()
	if err != nil {
		return 0, err
	}
	res, err := conn.Do("ZINCRBY", key, strconv.FormatFloat(score, 'f', -1, 64), value)
	if err != nil {
		return 0, err
	}
	data := string(res.([]byte))
	f, _ := strconv.ParseFloat(data, 64)
	return f, nil
}

func (r *myRedis) Zrank(key string, value string) (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return 0, err
	}
	res, err := conn.Do("ZRANK", key, value)
	if err != nil {
		return 0, err
	}
	return int(res.(int64)), nil
}

func (r *myRedis) Zrevrank(key string, value string) (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return 0, err
	}
	res, err := conn.Do("ZREVRANK", key, value)
	if err != nil {
		return 0, err
	}
	return int(res.(int64)), nil
}

func (r *myRedis) Zrange(key string, start int, end int) ([][]byte, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := conn.Do("ZRANGE", key, strconv.Itoa(start), strconv.Itoa(end))
	if err != nil {
		return nil, err
	}
	return res.([][]byte), nil
}

func (r *myRedis) Zrevrange(key string, start int, end int) ([][]byte, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := conn.Do("ZREVRANGE", key, strconv.Itoa(start), strconv.Itoa(end))
	if err != nil {
		return nil, err
	}
	return res.([][]byte), nil
}

func (r *myRedis) Zrangebyscore(key string, start float64, end float64) ([][]byte, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := conn.Do("ZRANGEBYSCORE", key, strconv.FormatFloat(start, 'f', -1, 64), strconv.FormatFloat(end, 'f', -1, 64))
	if err != nil {
		return nil, err
	}
	return res.([][]byte), nil
}

func (r *myRedis) Zcard(key string) (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("ZCARD", key)
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func (r *myRedis) Zscore(key string, member string) (float64, error) {
	conn, err := r.Dial()
	if err != nil {
		return 0, err
	}
	res, err := conn.Do("ZSCORE", key, member)
	if err != nil {
		return 0, err
	}
	data := string(res.([]byte))
	f, _ := strconv.ParseFloat(data, 64)
	return f, nil
}

func (r *myRedis) Zremrangebyrank(key string, start int, end int) (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("ZREMRANGEBYRANK", key, strconv.Itoa(start), strconv.Itoa(end))
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func (r *myRedis) Zremrangebyscore(key string, start float64, end float64) (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("ZREMRANGEBYSCORE", key, strconv.FormatFloat(start, 'f', -1, 64), strconv.FormatFloat(end, 'f', -1, 64))
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

// hash commands
func (r *myRedis) Hset(key string, field string, val string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("HSET", key, field, val)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Hmset(key string, params map[string](string)) error {
	conn, err := r.Dial()
	if err != nil {
		return err
	}
	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(params)...)
	if err != nil {
		return err
	}
	return nil
}

func (r *myRedis) Hget(key string, field string) (string, error) {
	conn, err := r.Dial()
	if err != nil {
		return "", err
	}

	res, err := redis.Bytes(conn.Do("HGET", key, field))
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}
	return string(res), nil
}

func (r *myRedis) Mget(keys []string) (map[string](string), error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(conn.Do("Mget", redis.Args{}.AddFlat(keys)...))
	if err != nil {
		return nil, err
	}
	var data = make(map[string](string), len(keys))
	for k, v := range keys {
		data[v] = res[k]
	}
	return data, nil
}

func (r *myRedis) Hmget(key string, params []string) (map[string](string), error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(conn.Do("Hmget", redis.Args{}.Add(key).AddFlat(params)...))
	if err != nil {
		return nil, err
	}
	var data = make(map[string](string), len(params))
	for k, v := range params {
		if res[k] != "" {
			data[v] = res[k]
		}
	}
	return data, nil
}

func (r *myRedis) Hgetall(key string) (map[string](string), error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	lenghth := len(res)
	var data = make(map[string](string), lenghth/2)
	for i := 0; i < lenghth; i = i + 2 {
		if res[i+1] != "" {
			data[res[i]] = res[i+1]
		}
	}
	return data, nil
}

func (r *myRedis) Hincrby(key string, field string, val int64) (int64, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("HINCRBY", key, field, strconv.FormatInt(val, 10))
	if err != nil {
		return -1, err
	}
	return res.(int64), nil
}

func (r *myRedis) Hexists(key string, field string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}

	res, err := conn.Do("HEXISTS", key, field)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Hdel(key string, field string) (bool, error) {
	conn, err := r.Dial()
	if err != nil {
		return false, err
	}
	res, err := conn.Do("HDEL", key, field)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *myRedis) Hlen(key string) (int, error) {
	conn, err := r.Dial()
	if err != nil {
		return -1, err
	}
	res, err := conn.Do("HLEN", key)
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func (r *myRedis) Hkeys(key string) ([]string, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *myRedis) Hvals(key string) ([]string, error) {
	conn, err := r.Dial()
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(conn.Do("HVALS", key))
	if err != nil {
		return nil, err
	}
	return res, nil
}
