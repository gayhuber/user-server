package hprose

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"testing"
)

type args map[string]interface{}

type RespBody struct {
	ErrorCode    int                    `json:"errorCode"`
	ErrorMsg     string                 `json:"errorMsg"`
	ResponseData map[string]interface{} `json:"responseData"`
}

func respHandler(res interface{}) (tmp map[string]interface{}) {
	// map 需要初始化一个出来
	tmp = make(map[string]interface{})
	log.Println("input res is : ", res)
	switch res.(type) {
	case nil:
		log.Printf("nil res: %v", res)
		return tmp
	case map[string]interface{}:
		log.Printf("map[string]interface{} res: %v", res)
		return res.(map[string]interface{})
	case map[interface{}]interface{}:
		log.Println("map[interface{}]interface{} res:", res)
		for k, v := range res.(map[interface{}]interface{}) {
			log.Printf("loop: k: %v, v: %v \n", k, v)
			switch k.(type) {
			case string:
				switch v.(type) {
				case map[interface{}]interface{}:
					log.Println("map[interface{}]interface{} v:", v)
					tmp[k.(string)] = respHandler(v)
					continue
				default:
					log.Printf("default value k: %v , v: %v \n", k, v)
					tmp[k.(string)] = v
					continue
				}

			default:
				log.Printf("default key k: %v , v: %v \n", k, v)
				continue
			}
		}
		return tmp
	default:
		// 暂时没遇到更复杂的数据
		log.Println("unknow data:", res)
	}
	return tmp
}

func map2Struct(mp map[string]interface{}, v interface{}) error {
	by, err := json.Marshal(mp)
	if err != nil {
		return nil
	}
	return json.Unmarshal(by, v)
}

func TestCallback(t *testing.T) {
	args := map[string]interface{}{
		"mobile": "18333636949",
	}

	// 注册方法
	ServiceMp.AddMethod("XyToken", BaseClient{
		Module: "System",
		Class:  "XyToken",
		Func:   "getXyOpenKey",
	})

	// 调用方法, 并传递参数进去, 获取远程 base 的调用结果
	res, err := RemoteFunc("XyToken", args)

	t.Log(res, err)
}

//
func TestCallback2(t *testing.T) {
	ServiceMp.AddMethod("getSimpleUserInfoById", BaseClient{
		Module: "User",
		Class:  "UserNew",
		Func:   "getSimpleUserInfoById",
	})

	field := []string{"uid", "user_name", "avatar", "login_mobile"}
	resp, err := GetSimpleUserInfoByID(20532239, false, field)

	if err != nil {
		t.Error(err)
	}

	fmt.Println("result:", resp)
}

func TestCallback3(t *testing.T) {
	ServiceMp.AddMethod("getXyOpenKey", BaseClient{
		Module: "System",
		Class:  "XyToken",
		Func:   "getXyOpenKey",
	})
	ServiceMp.AddMethod("quickMobileLogin", BaseClient{
		Module: "Passport",
		Class:  "Core\\Login",
		Func:   "quickMobileLogin",
	})

	code := "464101"
	mobile := "18610341055"

	resp, err := QuickMobileLogin(mobile, code, "0086", 11)

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("result:", resp)
	by, err := json.Marshal(resp)
	t.Logf("error: %v,json output:%s \n", err, string(by))
}

func TestCallback4(t *testing.T) {
	ServiceMp.AddMethod("getXyOpenKey", BaseClient{
		Module: "System",
		Class:  "XyToken",
		Func:   "getXyOpenKey",
	})

	ServiceMp.AddMethod("sendMobileSmsCode", BaseClient{
		Module: "Passport",
		Class:  "Core\\MobileSmsCode",
		Func:   "sendMobileSmsCode",
	})

	mobile := "18610341055"
	t.Logf("send message to %s", mobile)

	err := SendMobileSmsCode(mobile, "0086", "127.0.0.1", "8")

	if err != nil {
		t.Error(err)
	}
}

// SendMobileSmsCode 发送远程短信
func SendMobileSmsCode(mobile, countryCode, ip, Type string) error {
	openKey, err := GetXyOpenKey(mobile)
	if err != nil {
		return err
	}

	arg := args{
		"mobile":      mobile,
		"countryCode": countryCode,
		"type":        Type,
		"key":         openKey,
		"ip":          ip,
	}

	res, err := RemoteFunc("sendMobileSmsCode", arg)

	if err != nil {
		return err
	}

	tmp := respHandler(res)

	if tmp["errorCode"] != "0" {
		return errors.New(tmp["errorMsg"].(string))
	}
	return nil
}

// GetXyOpenKey 获取 openkey
func GetXyOpenKey(mobile string) (string, error) {
	res, err := RemoteFunc("getXyOpenKey", args{
		"mobile": mobile,
	})

	fmt.Println("remote openKey:", res, err)

	return res.(string), err
}

// GetSimpleUserInfoByID 获取用户信息
func GetSimpleUserInfoByID(uid int, returnUIDAsKey bool, arrField []string) (info map[string]interface{}, err error) {

	arg := args{
		"uid":            uid,
		"returnUidAsKey": returnUIDAsKey,
		"arrField":       arrField,
	}

	res, err := RemoteFunc("getSimpleUserInfoById", arg)
	fmt.Println("remote GetSimpleUserInfoById:", res, err)

	if err != nil {
		return
	}

	tmp := respHandler(res)
	fmt.Println("this is tmp:", tmp)
	return tmp, nil
}

// QuickMobileLogin 手机验证码快速注册/登录
func QuickMobileLogin(mobile, smsCode, countryCode string, sys int) (responseData map[string]interface{}, err error) {
	openKey, err := GetXyOpenKey(mobile)
	if err != nil {
		return
	}

	arg := args{
		"mobile":      mobile,
		"smsCode":     smsCode,
		"countryCode": countryCode,
		"key":         openKey,
		"extInfo": args{
			"sys": sys,
		},
	}

	res, err := RemoteFunc("quickMobileLogin", arg)
	fmt.Println("remote QuickMobileLogin:", res, err)

	if err != nil {
		return
	}

	tmp := respHandler(res)
	resp := RespBody{}
	err = map2Struct(tmp, &resp)

	// 每个接口对成功返回的定义还不一样....
	if resp.ErrorCode == 200 {
		return resp.ResponseData, nil
	}
	err = errors.New(resp.ErrorMsg)
	return
}
