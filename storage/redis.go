package storage

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
	e "webapp/apperrors"
	"webapp/stat"

	"github.com/zanlichard/redisgoe/redis"
)

const (
	StatCacheGet     = "RedisGet"
	StatCacheSet     = "RedisSet"
	StatCacheSete    = "RedisSetE"
	StatCacheDel     = "RedisDel"
	StatCacheTTL     = "RedisTTL"
	StatCachePublish = "RedisPublish"
	StatCacheMGet    = "RedisMGet"
	StatCacheMSet    = "RedisMSet"
	StatCacheLPush   = "RedisLPush"
)

var (
	defaultRedisTimeout = 30 * time.Second
	redisCachePool      *redis.Pool
	cacheRedis          = "cache"
	openStat            = false
)

func InitCache(redisHost string, redisPort int, redisPasswd string, maxIdleConn int, maxActiveConn int, idleTimeout int, statSwitch bool) *e.AppError {
	//Logger.Debug("InitCache host:%s port:%d", RedisHost, RedisPort)
	redisServer := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisCachePool = redisConnectPool(redisServer, redisPasswd, maxIdleConn, maxActiveConn, idleTimeout)
	if redisCachePool == nil {
		retCode := e.RetCode_ERR_CACHE_INIT
		return &e.AppError{Msg: e.RetCodeMsg[retCode], Err: nil, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	} else {
		if _, err := redisCachePool.Get().Do("PING"); err != nil {
			retCode := e.RetCode_ERR_CACHE_PING
			return &e.AppError{Msg: e.RetCodeMsg[retCode], Err: nil, Code: retCode, ErrPoint: e.GetErrPoint(1)}
		}
	}
	if statSwitch {
		openStat = true
		redisStat()
	}
	return nil
}

func redisStat() {
	stat.GStat.AddReportBodyRowItem(StatCacheGet)
	stat.GStat.AddReportBodyRowItem(StatCacheSet)
	stat.GStat.AddReportBodyRowItem(StatCacheDel)
	stat.GStat.AddReportBodyRowItem(StatCacheTTL)
	stat.GStat.AddReportBodyRowItem(StatCacheSete)
	stat.GStat.AddReportBodyRowItem(StatCacheMGet)
	stat.GStat.AddReportBodyRowItem(StatCachePublish)
	stat.GStat.AddReportBodyRowItem(StatCacheMSet)
	stat.GStat.AddReportBodyRowItem(StatCacheLPush)

	stat.GStat.AddReportErrorItem(StatCacheSete)
	stat.GStat.AddReportErrorItem(StatCacheGet)
	stat.GStat.AddReportErrorItem(StatCacheSet)
	stat.GStat.AddReportErrorItem(StatCacheDel)
	stat.GStat.AddReportErrorItem(StatCacheTTL)
	stat.GStat.AddReportErrorItem(StatCacheMGet)
	stat.GStat.AddReportErrorItem(StatCachePublish)
	stat.GStat.AddReportErrorItem(StatCacheMSet)
	stat.GStat.AddReportErrorItem(StatCacheLPush)

}

func redisConnectPool(server, passwd string, maxIdle int, maxActive int, idleTimeout int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		MaxActive:   maxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", server,
				defaultRedisTimeout,
				defaultRedisTimeout,
				defaultRedisTimeout)
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

func GetRedisConn(which string, timeleft int64) redis.Conn {
	if timeleft <= 0 {
		return nil
	}
	conn := redisCachePool.Get()
	if conn.GetConn() == nil {
		return nil
	}
	return conn

}

func RedisDel(conn redis.Conn, k string) *e.CallStack {
	st := e.BeginCallStack("redis.del")
	defer st.EndCall(1)
	retCode := e.RetCode_SUCCESS
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	_, err := conn.Do("DEL", k)
	if err != nil {
		retCode = (e.RetCode_ERR_CACHE_DEL)
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	}
	if openStat {
		stat.PushStat(StatCacheDel, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	return st
}

func RedisLPush(conn redis.Conn, k string, v string) (st *e.CallStack) {
	st = e.BeginCallStack("redis.lpush")
	defer st.EndCall(1)
	retCode := e.RetCode_SUCCESS
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	_, err := conn.Do("LPUSH", k, v)
	if err != nil {
		retCode = e.RetCode_ERR_CACHE_LPUSH
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	}
	if openStat {
		stat.PushStat(StatCacheLPush, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	return

}

func RedisMSet(conn redis.Conn, ks []string, datas []string, timeout int64) (st *e.CallStack) {
	st = e.BeginCallStack("redis.mset")
	defer st.EndCall(1)
	retCode := e.RetCode_SUCCESS
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])

	it := make([]interface{}, 2*len(ks), 2*len(ks))
	for i := 0; i < len(ks); i++ {
		it[2*i] = ks[i]
		it[2*i+1] = datas[i]
	}
	_, err := conn.Do("MSET", it...)
	if err != nil {
		retCode = e.RetCode_ERR_CACHE_MSET
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	}
	if openStat {
		stat.PushStat(StatCacheMSet, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	return
}

func RedisMGet(conn redis.Conn, ks []string, timeout int64) ([]string, *e.CallStack) {
	st := e.BeginCallStack("redis.mget")
	defer st.EndCall(1)
	var rs []string = nil
	retCode := e.RetCode_SUCCESS
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	mk := make([]interface{}, len(ks), len(ks))
	for i, v := range ks {
		mk[i] = v
	}

	ret, err := conn.Do("MGET", mk...)
	if err != nil || ret == nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		retCode = e.RetCode_ERR_CACHE_MGET
		goto RedisMGetStat
	}

	rs, err = redis.Strings(ret, err)
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		retCode = e.RetCode_ERR_TYPE_ASSERT
	}
RedisMGetStat:
	if openStat {
		stat.PushStat(StatCacheGet, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	return rs, st
}

func RedisGet(conn redis.Conn, k string, pValue interface{}, timeout int64) (bool, *e.CallStack) {
	retCode := e.RetCode_SUCCESS
	st := e.BeginCallStack("redis.get")
	defer st.EndCall(1)
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	ret, err := conn.Do("GET", k)
	if err != nil {
		retCode = (e.RetCode_ERR_CACHE_GET)
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	} else if ret == nil {
		retCode = (e.RetCode_ERR_CACHE_MISS)
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	} else {
		str, err := redis.String(ret, err)
		if err != nil {
			retCode = (e.RetCode_ERR_TYPE_ASSERT)
			st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
		} else {
			err = json.Unmarshal([]byte(str), &pValue)
			if err != nil {
				retCode = (e.RetCode_ERR_MARSH)
				st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}

			}
		}
	}
	if openStat {
		stat.PushStat(StatCacheGet, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	if retCode != 0 {
		return false, st
	}
	return true, st
}

func RedisSet(conn redis.Conn, k string, pValue interface{}, timeout int64) *e.CallStack {
	retCode := e.RetCode_SUCCESS
	st := e.BeginCallStack("redis.setex")
	defer st.EndCall(1)
	t1 := time.Now()
	b, err := json.Marshal(pValue)
	if err != nil {
		retCode = e.RetCode_ERR_MARSH
		st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
	} else {
		_, err = conn.Do("SET", k, string(b))
		if err != nil {
			retCode = e.RetCode_ERR_CACHE_SET
			st.ErrRet = &e.AppError{Err: err, Code: retCode, ErrPoint: e.GetErrPoint(1)}
		}
	}
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	if openStat {
		stat.PushStat(StatCacheSet, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	return st
}

// expire 单位为秒
func RedisSetEx(conn redis.Conn, k string, pValue interface{}, expire int, timeout int64) *e.CallStack {
	st := e.BeginCallStack("redis.setex")
	defer st.EndCall(1)
	retCode := e.RetCode_SUCCESS
	t1 := time.Now()
	b, err := json.Marshal(pValue)
	if err != nil {
		retCode = e.RetCode_ERR_MARSH
		st.ErrRet = &e.AppError{Code: retCode, Err: err, ErrPoint: e.GetErrPoint(1)}
	} else {
		_, err = conn.Do("SETEX", k, expire, string(b))
		if err != nil {
			retCode = e.RetCode_ERR_CACHE_SETE
			st.ErrRet = &e.AppError{Code: retCode, Err: err, ErrPoint: e.GetErrPoint(1)}
		}
	}
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	if openStat {
		stat.PushStat(StatCacheSete, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	return st
}
func RedisTTL(conn redis.Conn, k string, timeout int64) (int64, *e.CallStack) {
	st := e.BeginCallStack("redis.ttl")
	defer st.EndCall(1)
	retCode := e.RetCode_SUCCESS
	t1 := time.Now()
	timeLeft := int64(0)
	ttlTime, err := conn.Do("TTL", k)
	if err != nil {
		retCode = e.RetCode_ERR_CACHE_TTL
		st.ErrRet = &e.AppError{Code: retCode, Err: err, ErrPoint: e.GetErrPoint(1)}
		timeLeft = int64(e.RetCode_ERR_CACHE_TTL)

	} else {
		timeLeft2, err2 := redis.Int64(ttlTime, err)
		timeLeft = timeLeft2
		if err2 != nil {
			retCode = e.RetCode_ERR_TYPE_ASSERT
			st.ErrRet = &e.AppError{Code: retCode, Err: err2, ErrPoint: e.GetErrPoint(1)}
		}

	}
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	if openStat {
		stat.PushStat(StatCacheTTL, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}
	return timeLeft, st

}

func RedisMSetEx(conn redis.Conn, ks []string, datas []string, expire int, timeout int64) (st *e.CallStack) {
	st = e.BeginCallStack("redis.msetex")
	defer st.EndCall(1)
	retCode := e.RetCode_SUCCESS
	t1 := time.Now()
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])

	for i := range ks {
		if err := conn.Send("SETEX", ks[i], expire, datas[i]); err != nil {
			st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
			retCode = e.RetCode_ERR_CACHE_SETX
			goto RedisMSetExStat
		}
	}

	if err := conn.Flush(); err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		retCode = e.RetCode_ERR_CACHE_FLUSH
	}

RedisMSetExStat:
	if openStat {
		stat.PushStat(StatCacheTTL, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}

	return
}

func RedisPublish(conn redis.Conn, channel, pValue interface{}, timeout int64) *e.CallStack {
	retCode := e.RetCode_SUCCESS
	st := e.BeginCallStack("redis.publish")
	defer st.EndCall(1)
	srcAddr := net.ParseIP(strings.Split(conn.GetConn().RemoteAddr().String(), ":")[0])
	t1 := time.Now()
	b, err := json.Marshal(pValue)
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		retCode = e.RetCode_ERR_MARSH
		goto redisPubStat
	}
	_, err = conn.Do("PUBLISH", channel, string(b))
	if err != nil {
		st.ErrRet = &e.AppError{Err: err, ErrPoint: e.GetErrPoint(1)}
		retCode = e.RetCode_ERR_CACHE_PUB
	}
redisPubStat:
	if openStat {
		stat.PushStat(StatCachePublish, int(time.Now().Sub(t1).Seconds()*1000), srcAddr, 0, int(retCode))
	}

	return st
}
