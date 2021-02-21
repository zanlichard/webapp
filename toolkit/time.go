package toolkit

import (
	"time"
)

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GetTimeStamp() int64 {
	return time.Now().UTC().Unix()
}

func GetMSTimeStamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func GetSecTimeStamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Second)
}

func Time2MSTimeStamp(t *time.Time) int64 {
	ts := t.UTC().UnixNano() / int64(time.Millisecond)
	if ts > 0 {
		return ts
	} else {
		return 0
	}
}

type CustomTime struct {
	time.Time
}

const (
	ctLayout = "2006/01/02|15:04:05"
)

var nilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	ct.Time, err = time.Parse(ctLayout, string(b))
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format(ctLayout)), nil
}

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}
