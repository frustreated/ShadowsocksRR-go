// ssr_local
package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"

	ssr "github.com/10bits/shadowsocksR"
	"github.com/10bits/shadowsocksRR-go/utils"
)

var server *utils.SSRServer

func init() {
	room_url := "ssr://MTk0LjE1Ni4yMzAuMTk5OjMxNTYzOmF1dGhfYWVzMTI4X21kNTpzYWxzYTIwOmh0dHBfcG9zdDpWakZuVTJKaFMydDZRZy8_b2Jmc3BhcmFtPVpuUndMbXB3TG1SbFltbGhiaTV2Y21jJmdyb3VwPWFIUjBjRG92TDNkM2R5NXpjM0oxTG1sdVptOCZyZW1hcmtzPTQ0Q1E1WVdONkxTNTZJcUM1NEs1NDRDUjVweXQ1Ym1N"
	server = utils.SSRDecodeUrl(room_url)
}

func handleConnection(local net.Conn) {
	log.Println("handle connection")
	closed := false
	var err error
	defer func() {
		if !closed {
			local.Close()
		}
	}()
	r := bufio.NewReader(local)
	err = utils.Socks5HandShake(r, local)
	if err != nil {
		log.Println(err)
		return
	}
	addr_s, err := utils.Socks5ReadAddr(r)
	if err != nil {
		log.Println(err)
		return
	}
	err = utils.Socks5ConnectionConfirm(local)
	if err != nil {
		log.Println(err)
		return
	}
	//rawaddr := socks.ParseAddr(addr_s)
	rawaddr := addr_s
	const maxTryCount = 3
	var remote net.Conn
	for i := 0; i < maxTryCount; i++ {
		remote, err = connectToServer(rawaddr)
		if err != nil {
			log.Println("reconnect remote server")
		} else {
			break
		}
	}
	if err != nil {
		log.Println("cannot connect remote server")
		return
	}
	log.Println("connect remote server ok")
	defer func() {
		if !closed {
			remote.Close()
		}
	}()
	go gophplib.PipeThenClose(local, remote, nil)
	gophplib.PipeThenClose(remote, local, nil)
	closed = true
	log.Println("close connection.")
}

func connectToServer(rawaddr socks.Addr) (net.Conn, error) {
	u := utils.NewSSRUrlQuery(server)
	ssrconn, err := ssr.NewSSRClient(u)
	if err != nil {
		return nil, fmt.Errorf("connecting to SSR server failed :%v", err)
	}
	if server.ObfsData == nil {
		server.ObfsData = ssrconn.IObfs.GetData()
	}
	ssrconn.IObfs.SetData(server.ObfsData)

	if server.ProtocolData == nil {
		server.ProtocolData = ssrconn.IProtocol.GetData()
	}
	ssrconn.IProtocol.SetData(server.ProtocolData)

	if _, err := ssrconn.Write(rawaddr); err != nil {
		ssrconn.Close()
		return nil, err
	}
	return ssrconn, nil

}
func runLocal() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

// func main() {
// 	runLocal()
// }
