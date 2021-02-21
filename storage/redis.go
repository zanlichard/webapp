package storage

import (
	"encoding/json"
	"net"
	"strings"
	"time"
	e "webapp/errors"
	"webapp/stat"
	"github.com/garyburd/redigo/redis"
)

const (
	StatCacheGet  = "RedisGet"
	StatCacheSet  = "RedisSet"
	StatCacheSete = "RedisSetE"
	StatCacheDel  = "RedisDel"
	StatCacheTTL  = "RedisTTL"
)

var (
	DefaultRedisTimeout = 30 * time.Second
)

func InitRedisStat() {
	stat.GStat.AddReportBodyRowItem(StatCacheGet)
	stat.GStat.AddReportBodyRowItem(StatCacheSet)
	stat.GStat.AddReportBodyRowItem(StatCacheDel)
	stat.GStat.AddReportBodyRowItem(StatCacheTTL)
	stat.GStat.AddReportBodyRowItem(StatCacheSete)
	stat.GStat.AddReportErrorItem(StatCacheSete)
	stat.GStat.AddReportErrorItem(StatCacheGet)
	stat.GStat.AddReportErrorItem(StatCacheSet)
	stat.GStat.AddReportErrorItem(StatCacheDel)
	stat.GStat.AddReportErrorItem(StatCacheTTL)
}
func RedisDel(conn redis.Conn, k string) *e.CallStack {
	st := e.BeginCallStack("redis.del")
	defer st.EndCall(1)
	retCode := 0
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	_, err := conn.Do("DEL", k)
	if err != nil {
		retCode = int(e.RetCode_ERR_CACHE_DEL)
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	}
	stat.PushStat(StatCacheDel, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, retCode)
	return st
}

func RedisLPush(conn redis.Conn, k string, v string) (st *e.CallStack) {
	st = e.BeginCallStack("redis.lpush")
	defer st.EndCall(1)

	_, err := conn.Do("LPUSH", k, v)
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		return
	}

	return
}

func RedisMSet(conn redis.Conn, ks []string, datas []string, timeout int64) (st *e.CallStack) {
	st = e.BeginCallStack("redis.mset")
	defer st.EndCall(1)

	it := make([]interface{}, 2*len(ks), 2*len(ks))
	for i := 0; i < len(ks); i++ {
		it[2*i] = ks[i]
		it[2*i+1] = datas[i]
	}

	_, err := conn.Do("MSET", it...)
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		return
	}

	return
}

func RedisMGet(conn redis.Conn, ks []string, timeout int64) (rs []string, st *e.CallStack) {
	st = e.BeginCallStack("redis.mget")
	defer st.EndCall(1)

	mk := make([]interface{}, len(ks), len(ks))
	for i, v := range ks {
		mk[i] = v
	}

	ret, err := conn.Do("MGET", mk...)
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		return nil, st
	}
	if ret == nil {
		return nil, st
	}

	rs, err = redis.Strings(ret, err)
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		return nil, st
	}

	return rs, st
}

func RedisGet(conn redis.Conn, k string, pValue interface{}, timeout int64) (bool, *e.CallStack) {
	retCode := 0
	st := e.BeginCallStack("redis.get")
	defer st.EndCall(1)
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	ret, err := conn.Do("GET", k)
	if err != nil {
		retCode = int(e.RetCode_ERR_CACHE_GET)
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	} else if ret == nil {
		retCode = int(e.RetCode_ERR_CACHE_MISS)
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	} else {
		str, err := redis.String(ret, err)
		if err != nil {
			retCode = int(e.RetCode_ERR_TYPE_ASSERT)
			st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
		} else {
			err = json.Unmarshal([]byte(str), &pValue)
			if err != nil {
				retCode = int(e.RetCode_ERR_MARSH)
				st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}

			}
		}
	}
	stat.PushStat(StatCacheGet, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, retCode)
	if retCode != 0 {
		return false, st
	}
	return true, st
}

func RedisSet(conn redis.Conn, k string, pValue interface{}, timeout int64) *e.CallStack {
	retCode := 0
	st := e.BeginCallStack("redis.setex")
	defer st.EndCall(1)
	t1 := time.Now()
	b, err := json.Marshal(pValue)
	if err != nil {
		retCode = int(e.RetCode_ERR_MARSH)
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	} else {
		_, err = conn.Do("SET", k, string(b))
		if err != nil {
			retCode = int(e.RetCode_ERR_CACHE_SET)
			st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
		}
	}
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	stat.PushStat(StatCacheSet, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, retCode)
	return st
}

// expire 单位为秒
func RedisSetEx(conn redis.Conn, k string, pValue interface{}, expire int, timeout int64) *e.CallStack {
	st := e.BeginCallStack("redis.setex")
	defer st.EndCall(1)
	retCode := 0
	t1 := time.Now()
	b, err := json.Marshal(pValue)
	if err != nil {
		retCode = int(e.RetCode_ERR_MARSH)
		st.ErrRet = &e.AppError{Code: retCode, Err: err, ErrPoint: e.GetErrPoint(1)}
	} else {
		_, err = conn.Do("SETEX", k, expire, string(b))
		if err != nil {
			retCode = int(e.RetCode_ERR_CACHE_SETE)
			st.ErrRet = &e.AppError{Code: retCode, Err: err, ErrPoint: e.GetErrPoint(1)}
		}
	}
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	stat.PushStat(StatCacheSete, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, retCode)
	return st
}
func RedisTTL(conn redis.Conn, k string, timeout int64) (int64, *e.CallStack) {
	st := e.BeginCallStack("redis.ttl")
	defer st.EndCall(1)
	retCode := 0
	t1 := time.Now()
	timeLeft := int64(0)
	ttlTime, err := conn.Do("TTL", k)
	if err != nil {
		retCode = int(e.RetCode_ERR_CACHE_TTL)
		st.ErrRet = &e.AppError{Code: retCode, Err: err, ErrPoint: e.GetErrPoint(1)}
		timeLeft = int64(e.RetCode_ERR_CACHE_TTL)

	} else {
		timeLeft2, err2 := redis.Int64(ttlTime, err)
		timeLeft = timeLeft2
		if err2 != nil {
			retCode = int(e.RetCode_ERR_TYPE_ASSERT)
			st.ErrRet = &e.AppError{Code: retCode, Err: err2, ErrPoint: e.GetErrPoint(1)}
			timeLeft = int64(e.RetCode_ERR_TYPE_ASSERT)
		}

	}
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	stat.PushStat(StatCacheTTL, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, retCode)
	return timeLeft, st

}

func RedisMSetEx(conn redis.Conn, ks []string, datas []string, expire int, timeout int64) (st *e.CallStack) {
	st = e.BeginCallStack("redis.msetex")
	defer st.EndCall(1)

	for i := range ks {
		if err := conn.Send("SETEX", ks[i], expire, datas[i]); err != nil {
			st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
			return
		}
	}

	if err := conn.Flush(); err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		return
	}

	return
}

func RedisPublish(conn redis.Conn, channel, pValue interface{}, timeout int64) *e.CallStack {
	st := e.BeginCallStack("redis.publish")
	defer st.EndCall(1)
	b, err := json.Marshal(pValue)
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		return st
	}
	_, err = conn.Do("PUBLISH", channel, string(b))
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		return st
	}
	return st
}

func RedisConnectPool(server, passwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", server,
				DefaultRedisTimeout,
				DefaultRedisTimeout,
				DefaultRedisTimeout)
			if err != nil {
				return nil, err
			}
			if passwd != "" {
				if _, err := c.Do("AUTH", passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
