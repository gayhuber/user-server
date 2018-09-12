package tools

import (
	"testing"
)

func TestGenerageUniqueID(t *testing.T) {

	uid, salt := GenerageUniqueID("hers")
	t.Logf("获取内容: uid: %s, salt: %s ", uid, salt)
}
