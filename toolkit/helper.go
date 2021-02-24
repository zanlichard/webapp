package toolkit

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type Int64Arr []int64

func (a Int64Arr) Len() int           { return len(a) }
func (a Int64Arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Int64Arr) Less(i, j int) bool { return a[i] < a[j] }

func RandomLowerLetterString(l int) string {
	var result bytes.Buffer
	var temp rune = 'a'
	for i := 0; i < l; {
		randX := RandomInt(97, 122)
		if rune(randX) != temp {
			temp = rune(randX)
			strChar := string(temp)
			result.WriteString(strChar)
			i++
		}
	}
	return result.String()
}

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func Md5Sum(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}

func ConvertToString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	}

	return ""
}
