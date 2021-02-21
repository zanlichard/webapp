package main

import "log"

func main() {
	token,err := NewJWT().CreateToken(1000009,"15532",30*24*3600)
	if err != nil {
		log.Printf("createtoken err: %v\n",err)
		return
	}
	log.Printf("token ==%v",token)
}
