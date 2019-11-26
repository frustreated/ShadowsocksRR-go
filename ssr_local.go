// ssr_local
package main

import (
	"log"
	"net"

	"github.com/10bits/shadowsocksR"
	"github.com/10bits/shadowsocksRR-go/utils"
)

var local_addr = "127.0.0.1:8081"
var server_addr = "127.0.0.1:8082"

func ssrServer() {
	l, err := net.Listen("tcp", server_addr)
	log.Println("server listening on:", server_addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go utils.HandleSocks5Conn(conn)
	}
}
func ssrLocal() {
	l, err := net.Listen("tcp", local_addr)
	log.Println("loacl listening on:", local_addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go utils.HandleLocalConn(conn, server_addr)
	}
}
func main() {
	go ssrServer()
	ssrLocal()
}
