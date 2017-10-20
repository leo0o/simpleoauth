package simpleoauth

import (
	"fmt"
	"httplib"
)

const weibo_getaccesstoken_url  = "https://api.weibo.com/oauth2/access_token"
const weibo_getuserinfo_url = "https://api.weibo.com/2/users/show.json"

var weiboOAuth = &WeiboOAuth{}

type WeiboOAuth struct {
	appkey string
	appsecret string
	redirect_url string
}

func (oauth *WeiboOAuth) GetAccesstoken(code string) map[string]interface{}{
	request:= httplib.Post(weibo_getaccesstoken_url)
	request.Param("client_id", oauth.appkey)
	request.Param("client_secret", oauth.appsecret)
	request.Param("grant_type", "authorization_code")
	request.Param("code", code)
	request.Param("redirect_uri", oauth.redirect_url)
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func (oauth *WeiboOAuth) GetUserinfo(accesstoken string, openid string) map[string]interface{}{
	request:= httplib.Get(weibo_getuserinfo_url)
	request.Param("access_token", accesstoken)
	request.Param("uid", openid)
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func (oauth *WeiboOAuth) Authorize(code string) AuthorizeResult{
	accesstokenResponse := oauth.GetAccesstoken(code)
	if accesstokenResponse == nil{
		return AuthorizeResult{false, nil}
	}
	_, ok := accesstokenResponse["error_code"]         //获取accesstoken接口返回错误码
	if ok {
		return AuthorizeResult{false, nil}
	}
	openid := accesstokenResponse["uid"].(string)
	accesstoken := accesstokenResponse["access_token"].(string)
	getuserinfoResult := oauth.GetUserinfo(accesstoken, openid)
	fmt.Println(getuserinfoResult)
	if getuserinfoResult == nil {
		return AuthorizeResult{false, nil}
	}
	_, ok = getuserinfoResult["error_code"]           //获取用户信息接口返回错误码
	if ok {
		return AuthorizeResult{false, nil}
	}
	var sex int
	if getuserinfoResult["gender"].(string) == "m"{
		sex = 1
	}else if getuserinfoResult["gender"].(string) == "f"{
		sex = 2
	}else if getuserinfoResult["gender"].(string) == "n"{
		sex = 0
	}
	return AuthorizeResult{true, map[string]interface{}{
		"nickname":getuserinfoResult["screen_name"].(string),
		"openid":openid,
		"sex":sex,
		"headimgurl":getuserinfoResult["profile_image_url"].(string),
		"unionid":""}}
}

func (oauth *WeiboOAuth) InitOAuth(){
	oauth.appkey = Weiboappkey
	oauth.appsecret = Weiboappsecret
	oauth.redirect_url = WeiboRedirectUrl
}


func init(){
	ReisterPlatform("weibo", weiboOAuth)
}

