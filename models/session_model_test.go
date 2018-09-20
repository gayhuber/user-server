package models

import (
	"testing"
)

type H map[string]interface{}

func TestStoreSession(t *testing.T) {
	info := H{
		"open_id": "a2a666556c8e804d9cce1266692ac76d",
		"sy_uid":  23282464,
		"src":     "hers",
		"info": H{
			"avatar":   "http://img2.soyoung.com/user/2.png",
			"nickname": "氧气_hers_1809183427",
		},
		"token": "74322a7baa1b4aedb3d704b21e76e843",
	}

	StoreSession("session_74322a7baa1b4aedb3d704b21e76e843", info)
}

func TestNewSession(t *testing.T) {

	token := "67bdd8840c993af18650423e219b1238"
	tp := "h5"
	resp, err := NewSession(token, tp).info()
	t.Log(string(resp), err)
}
