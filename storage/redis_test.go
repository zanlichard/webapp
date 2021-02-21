package storage

import (
	"encoding/json"
	"testing"
)

/*
type User struct {
	Uid    int64
	UserId string
	Email  string
}
*/

func TestRedisSet(t *testing.T) {

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

	user.Uid = 2
	user.UserId = "boa@china"
	st = SetCache("user2", &user, 1000)
	if st.ErrRet != nil {
		t.Error("redis cache set error:", st.GetProcMsg(1000))
		return
	}

	user.Uid = 3
	user.UserId = "boa@korea"
	st = SetCache("user3", &user, 1000)
	if st.ErrRet != nil {
		t.Error("redis cache set error:", st.GetProcMsg(1000))
		return
	}

}

func TestCacheMGet(t *testing.T) {

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

	t.Log("uid:", us[0].Uid, "userid:", us[0].UserId, " email:", us[0].Email)
	t.Log("uid:", us[1].Uid, "userid:", us[1].UserId, " email:", us[1].Email)
	RedisDel(conn, "user1")
	RedisDel(conn, "user2")
	RedisDel(conn, "user3")

}
