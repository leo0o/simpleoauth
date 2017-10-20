package simpleoauth

import (
	"httplib"
)

const wechat_getaccesstoken_url  = "https://api.weixin.qq.com/sns/oauth2/access_token"
const wechat_getuserinfo_url = "https://api.weixin.qq.com/sns/userinfo"

var wechatOAuth = &WechatOAuth{}

type WechatOAuth struct {
	appkey string
	appsecret string
}

func (oauth *WechatOAuth) GetAccesstoken(code string) map[string]interface{}{
	request:= httplib.Get(wechat_getaccesstoken_url)
	request.Param("appid",oauth.appkey)
	request.Param("secret",oauth.appsecret)
	request.Param("code",code)
	request.Param("grant_type","authorization_code")
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func (oauth *WechatOAuth) GetUserinfo(accesstoken string, openid string) map[string]interface{}{
	request:= httplib.Get(wechat_getuserinfo_url)
	request.Param("access_token", accesstoken)
	request.Param("openid", openid)
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func (oauth *WechatOAuth) Authorize(code string) AuthorizeResult{
	accesstokenResponse := oauth.GetAccesstoken(code)
	if accesstokenResponse == nil{
		return AuthorizeResult{false, nil}
	}
	_, ok := accesstokenResponse["errcode"]         //获取accesstoken接口返回错误码
	if ok {
		return AuthorizeResult{false, nil}
	}
	openid := accesstokenResponse["openid"].(string)
	accesstoken := accesstokenResponse["access_token"].(string)
	getuserinfoResult := oauth.GetUserinfo(accesstoken, openid)
	if getuserinfoResult == nil {
		return AuthorizeResult{false, nil}
	}
	_, ok = getuserinfoResult["errcode"]           //获取用户信息接口返回错误码
	if ok {
		return AuthorizeResult{false, nil}
	}

	return AuthorizeResult{true, map[string]interface{}{
		"nickname":getuserinfoResult["nickname"].(string),
		"openid":getuserinfoResult["openid"].(string),
		"sex":getuserinfoResult["sex"].(float64),
		"headimgurl":getuserinfoResult["headimgurl"].(string),
		"unionid":getuserinfoResult["unionid"].(string)}}
}

func (oauth *WechatOAuth) InitOAuth(){
	oauth.appkey = Wechatappkey
	oauth.appsecret = Wechatappsecret
}


func init(){
	ReisterPlatform("wechat", wechatOAuth)
}

