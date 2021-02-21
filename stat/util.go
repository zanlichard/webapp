package stat

import (
	"net"
	"strconv"
	"strings"
)

// Convert net.IP to int
func inet_aton(ipnr net.IP) int {
	if len(ipnr) <= 0 {
		return 0
	}
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int

	sum += int(b0) << 24
	sum += int(b1) << 16
	sum += int(b2) << 8
	sum += int(b3)

	return sum
}

func inet_ntoa(ipnr int) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func GetTopn(rawdata map[int]uint, topn int) (arraykeys []int) {
	arraykeys = make([]int, len(rawdata))
	for k, _ := range rawdata {
		arraykeys = append(arraykeys, k)
	}
	for i := 0; i < len(arraykeys); i++ {
		tempIndex := i
		for j := i + 1; j < len(arraykeys); j++ {
			if rawdata[arraykeys[j]] > rawdata[arraykeys[tempIndex]] {
				tempIndex = j
			}
		}
		if tempIndex != i {
			arraykeys[tempIndex], arraykeys[i] = arraykeys[i], arraykeys[tempIndex]
		}
	}
	return

}
