// ssr
package utils

import (
	"encoding/base64"
	"net"
	"net/url"
	"strconv"
	"strings"

	SSR "github.com/10bits/shadowsocksR"
	shadowsocksr "github.com/10bits/shadowsocksR"
	"github.com/10bits/shadowsocksR/obfs"
	"github.com/10bits/shadowsocksR/protocol"
	"github.com/10bits/shadowsocksR/ssr"
	"github.com/10bits/shadowsocksRR-go/enums"
)

type SSInfo struct {
	SSRInfo
	EncryptMethod   string
	EncryptPassword string
}
type SSRInfo struct {
	Obfs          string
	ObfsParam     string
	ObfsData      interface{}
	Protocol      string
	ProtocolParam string
	ProtocolData  interface{}
}
type SSRServer struct {
	SSInfo
	Address string
	Type    string
	Remarks string
	Group   string
}

func Base64Decode(s string) string {
	var b []byte
	b, _ = base64.RawURLEncoding.DecodeString(s)
	return string(b)
}

//ssr_url := "ssr://MTA4LjE3Ny4yMzUuMTI2OjM2MjUwOmF1dGhfYWVzMTI4X21kNTpjaGFjaGEyMC1pZXRmOmh0dHBfcG9zdDpNMnBqTURsd2NVOUlVZy8_b2Jmc3BhcmFtPVpuUndMblZ6TG1SbFltbGhiaTV2Y21jJnJlbWFya3M9NTc2TzVadTlRUSZncm91cD02WVc0NWFXMg"
func SSRDecodeUrl(ssr_url string) (ssrserver *SSRServer) {
	ssr_url = strings.TrimLeft(ssr_url, "ssr://")
	ssr_url = Base64Decode(ssr_url)
	v := strings.Split(ssr_url, ":")
	var ip, port, protocol, encrypt, obfs = v[0], v[1], v[2], v[3], v[4]
	u, _ := url.Parse("ssr://" + v[5])
	password := Base64Decode(u.Host)
	obfsparam := Base64Decode(u.Query().Get("obfsparam"))
	protocolparam := Base64Decode(u.Query().Get("protocolparam"))
	remarks := Base64Decode(u.Query().Get("remarks"))
	group := Base64Decode(u.Query().Get("group"))
	ssrserver = &SSRServer{
		Address: ip + ":" + port,
		Type:    "ssr",
		Remarks: remarks,
		Group:   group,
		SSInfo: SSInfo{
			EncryptMethod:   encrypt,
			EncryptPassword: password,
			SSRInfo: SSRInfo{
				Protocol:      protocol,
				ProtocolParam: protocolparam,
				Obfs:          obfs,
				ObfsParam:     obfsparam,
			},
		},
	}
	return
}

var DefaultServerAddr = "127.0.0.1:12345"
var DefaultSSRServer = &SSRServer{
	Address: DefaultServerAddr,
	Type:    "ssr",
	Remarks: "test",
	Group:   "test",
	SSInfo: SSInfo{
		EncryptMethod:   enums.ENUM_ENCRYPT_MAP[enums.ENCRYPT_CHACHA20_IETF],
		EncryptPassword: "test",
		SSRInfo: SSRInfo{
			Protocol:      enums.ENUM_PROTOCOL_MAP[enums.PROTOCOL_AUTH_AES128_MD5],
			ProtocolParam: "",
			Obfs:          enums.ENUM_OBFS_MAP[enums.OBFS_HTTP_POST],
			ObfsParam:     "",
		},
	},
}

func NewSSRUrlQuery(server *SSRServer) *url.URL {
	u := &url.URL{
		Scheme: server.Type,
		Host:   server.Address,
	}
	v := u.Query()
	v.Set("encrypt-method", server.EncryptMethod)
	v.Set("encrypt-key", server.EncryptPassword)
	v.Set("obfs", server.Obfs)
	v.Set("obfs-param", server.ObfsParam)
	v.Set("protocol", server.Protocol)
	v.Set("protocol-param", server.ProtocolParam)
	u.RawQuery = v.Encode()
	return u
}
func NewSSRTcpConn(u *url.URL, conn net.Conn) (*shadowsocksr.SSTCPConn, error) {
	var query = u.Query()
	encryptMethod := query.Get("encrypt-method")
	encryptKey := query.Get("encrypt-key")
	cipher, err := SSR.NewStreamCipher(encryptMethod, encryptKey)
	if err != nil {
		return nil, err
	}
	ssconn := SSR.NewSSTCPConn(conn, cipher)
	// should initialize obfs/protocol now
	addr_s := u.Host
	rs := strings.Split(addr_s, ":")
	port, _ := strconv.Atoi(rs[1])
	ssconn.IObfs = obfs.NewObfs(query.Get("obfs"))
	obfsServerInfo := &ssr.ServerInfoForObfs{
		Host:   rs[0],
		Port:   uint16(port),
		TcpMss: 1460,
		Param:  query.Get("obfs-param"),
	}
	ssconn.IObfs.SetServerInfo(obfsServerInfo)
	ssconn.IProtocol = protocol.NewProtocol(query.Get("protocol"))
	protocolServerInfo := &ssr.ServerInfoForObfs{
		Host:   rs[0],
		Port:   uint16(port),
		TcpMss: 1460,
		Param:  query.Get("protocol-param"),
	}
	ssconn.IProtocol.SetServerInfo(protocolServerInfo)

	return ssconn, nil
}
