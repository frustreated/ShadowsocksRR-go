package utils

import (
	"bufio"
	"log"
	"net"
)

func HandleSocks5Conn(local_conn net.Conn) {
	var closed bool = false
	var err error
	r := bufio.NewReader(local_conn)
	defer func() {
		if !closed {
			local_conn.Close()
		}
	}()
	err = Socks5HandShake(r, local_conn)
	if err != nil {
		log.Println("socks5 handshake:", err)
		return
	}
	addr_s, err := Socks5ReadAddr(r)
	if err != nil {
		log.Println("socks5 readaddr:", err)
		return
	}
	err = Socks5ConnectionConfirm(local_conn)
	if err != nil {
		log.Println("socks5 confirm:", err)
		return
	}
	log.Println(addr_s)
	var remote net.Conn
	remote, err = net.Dial("tcp", addr_s)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()
	go PipeThenClose(local_conn, remote, nil)
	PipeThenClose(remote, local_conn, nil)
	closed = true
}
func HandleLocalConn(local_conn net.Conn, remote_addr string) {
	closed := false
	defer func() {
		if !closed {
			local_conn.Close()
		}
	}()
	remote_conn, err := net.Dial("tcp", remote_addr)
	if err != nil {
		log.Println("connect remote server:", remote_addr, err)
		return
	}
	defer func() {
		if !closed {
			remote_conn.Close()
		}
	}()
	go PipeThenClose(local_conn, remote_conn, nil)
	PipeThenClose(remote_conn, local_conn, nil)
	closed = true
}
