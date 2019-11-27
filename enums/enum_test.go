package enums

import (
	"testing"
)

func TestEnums(t *testing.T) {
	t.Log(ENUM_ENCRYPT_MAP[ENCRYPT_AES_128_CTR])
	t.Log(ENUM_PROTOCOL_MAP[PROTOCOL_AUTH_AES128_MD5])
	t.Log(ENUM_OBFS_MAP[OBFS_PLAIN])
}
