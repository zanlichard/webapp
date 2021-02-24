package toolkit

import (
	"crypto/md5"
	"fmt"
	"sort"
)

func Md5Digest(data string) string {
	dataBytes := []byte(data)
	digestBytes := md5.Sum(dataBytes)
	md5Str := fmt.Sprintf("%x", digestBytes)
	return md5Str
}

func ApiSign(req string, key string) string {
	signStr := req + "_" + key
	return Md5Digest(signStr)
}

func ArrayCheckIn(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	//index的取值：[0,len(str_array)]
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}
