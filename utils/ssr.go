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
