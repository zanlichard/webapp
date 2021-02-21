package toolkit

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type NetIfs map[string]net.Interface

func GetNetIfs() (NetIfs, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ifs := make(NetIfs)
	// handle err
	for _, i := range ifaces {
		ifs[i.Name] = i
	}
	return ifs, nil
}

func GetNicIPs(ifname string) ([]net.IP, error) {
	ifs, err := GetNetIfs()
	if err != nil {
		return nil, err
	}
	nic, ok := ifs[ifname]
	if !ok {
		return nil, fmt.Errorf("cannot find ifname %s", ifname)
	}
	addrs, err := nic.Addrs()
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func IpIntToStringBig(ipInt int) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[i] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}
