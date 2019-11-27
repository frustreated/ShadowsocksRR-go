package utils

import (
	"testing"
)

func TestDecodeUrl(t *testing.T) {
	ssr_url := "ssr://MTA4LjE3Ny4yMzUuMTI2OjM2MjUwOmF1dGhfYWVzMTI4X21kNTpjaGFjaGEyMC1pZXRmOmh0dHBfcG9zdDpNMnBqTURsd2NVOUlVZy8_b2Jmc3BhcmFtPVpuUndMblZ6TG1SbFltbGhiaTV2Y21jJnJlbWFya3M9NTc2TzVadTlRUSZncm91cD02WVc0NWFXMg"
	server := SSRDecodeUrl(ssr_url)
	t.Log("地址:", server.Address)
	t.Log("加密方式:", server.EncryptMethod)
	t.Log("密码:", server.EncryptPassword)
	t.Log("协议:", server.Protocol)
	t.Log("协议参数:", server.ProtocolParam)
	t.Log("混淆:", server.Obfs)
	t.Log("混淆参数:", server.ObfsParam)
	t.Log("分组:", server.Group)
	t.Log("备注:", server.Remarks)
	if server.EncryptMethod != DefaultSSRServer.EncryptMethod {
		t.Error("test error")
	}
	if server.Protocol != DefaultSSRServer.Protocol {
		t.Error("test error")
	}
	if server.Obfs != DefaultSSRServer.Obfs {
		t.Error("test error")
	}

}
