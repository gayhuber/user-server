package models_test

import (
	"testing"
	"user-server/models"
)

func TestUniqueID(t *testing.T) {
	auth := models.H5Auth{Src: "hers"}

	t.Log("获取内容:", auth)
}
