package storage

import (
	"encoding/json"
	"testing"
)

type User struct {
	Uid    int64
	UserId string
	Email  string
}

func TestRedisSet(t *testing.T) {
	InitCache("192.168.163.129", 6379, "hsb_redis_123", false)
	conn, err := GetRedisConn(CacheRedis, 3000)
	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		t.Error("TestRedisSet get connection error")
	}

	user := new(User)
	user.UserId = "boa@japan"
	user.Uid = 1
	user.Email = "dream@foxmail.com"

	st := RedisSet(conn, "user1", &user, 1000)
	if st.ErrRet != nil {
		t.Error("redis cache set error:", st.GetProcMsg(1000))
		return
	}

	user2 := new(User)
	_, st = RedisGet(conn, "user1", &user2, 1000)
	if st.ErrRet != nil {
		t.Error("redis cache set error:", st.GetProcMsg(1000))
		return
	}
	t.Logf("get value:%+v", user2)

}

func TestCacheMGet(t *testing.T) {
	InitCache("192.168.163.129", 6379, "hsb_redis_123", false)
	us := make([]*User, 4, 4)
	ks := []string{"user1", "user2", "user4", "user3"}

	conn, err := GetRedisConn(CacheRedis, 3000)
	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		t.Error("TestRedisSet get connection error")
	}

	rs, st := RedisMGet(conn, ks, 3000)

	if st.ErrRet != nil {
		t.Error("redis mget error")
		return
	}

	if rs == nil {
		t.Error("redis mget not ok")
		return
	}

	for i, r := range rs {
		if r != "" {
			err = json.Unmarshal([]byte(r), &us[i])
			if err != nil {
				t.Log("redis mget transfer error")
				continue
			}

		}
	}
	t.Logf("user1:%+v", us[0])
	t.Logf("user2:%+v", us[1])
	t.Logf("user3:%+v", us[2])
	t.Logf("user4:%+v", us[3])
	RedisDel(conn, "user1")
	RedisDel(conn, "user2")
	RedisDel(conn, "user3")
	RedisDel(conn, "user4")

}
