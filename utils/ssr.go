// ssr
package utils

import (
	"encoding/base64"
	"net/url"
	"strings"
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
}

func NewServerFromUrl(room_url string) *SSRServer {
	var ssrserver *SSRServer
	room_url = strings.TrimLeft(room_url, "ssr://")
	var b []byte
	b, _ = base64.RawURLEncoding.DecodeString(room_url)
	room_url = string(b)
	room_url = "ssr://" + room_url
	u, _ := url.Parse(room_url)
	l := strings.Split(u.Host, ":")
	obfsparam := u.Query().Get("obfsparam")
	b, _ = base64.RawURLEncoding.DecodeString(obfsparam)
	obfsparam = string(b)
	var ip, port, password, protocol, obfs = l[0], l[1], l[5], l[2], l[4]
	b, _ = base64.RawURLEncoding.DecodeString(password)
	password = string(b)
	ssrserver = &SSRServer{
		Address: ip + ":" + port,
		Type:    "ssr",
		SSInfo: SSInfo{
			EncryptMethod:   l[3],
			EncryptPassword: password,
			SSRInfo: SSRInfo{
				Protocol:      protocol,
				ProtocolParam: "",
				Obfs:          obfs,
				ObfsParam:     obfsparam,
			},
		},
	}
	return ssrserver
}
