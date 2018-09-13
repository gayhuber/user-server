package tools

import (
	"testing"
)

func TestGenerageUniqueID(t *testing.T) {

	uid, salt := GenerageUniqueID("hers")
	t.Logf("获取内容: uid: %s, salt: %s ", uid, salt)
}

func TestRandInt(t *testing.T) {
	// randNum := RandInt(1000, 2000)
	t.Logf("获取随机数: %d ", RandInt(1000, 9999))
	t.Logf("获取随机数: %d ", RandInt(1, 5))
}
