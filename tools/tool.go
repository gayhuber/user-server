package tools

import (
	"crypto/md5"
	"fmt"
	"github.com/satori/go.uuid"
)

// GenerageUniqueID 获取uniqueid
func GenerageUniqueID(src string) (uid, salt string) {
	uuid, _ := uuid.NewV4()
	salt = uuid.String()[:32]
	h := md5.New()
	h.Write([]byte(src + salt))
	uid = fmt.Sprintf("%x\n", h.Sum(nil))
	return
}

// MD5String 生成 md5 字符串
func MD5String(v []byte) string {
	h := md5.New()
	h.Write(v)
	return fmt.Sprintf("%x", h.Sum(nil))
}
