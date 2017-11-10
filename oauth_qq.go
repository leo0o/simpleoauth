package simpleoauth

import (
	"httplib"
	"strings"
	"encoding/json"
)

const qq_getaccesstoken_url = "https://graph.qq.com/oauth2.0/token"
const qq_getuserinfo_url = "https://graph.qq.com/user/get_user_info"
const qq_openid_url = "https://graph.qq.com/oauth2.0/me"

var qqOAuth = &QQOAuth{}

type QQOAuth struct {
	appkey string
	appsecret string
	redirect_url string
}

func (oauth *QQOAuth) GetAccesstoken(code string) map[string]interface{}{
	request:= httplib.Get(qq_getaccesstoken_url)
	request.Param("grant_type","authorization_code")
	request.Param("client_id",oauth.appkey)
	request.Param("client_secret",oauth.appsecret)
	request.Param("code",code)
	request.Param("redirect_uri", oauth.redirect_url)

	response, err := request.String()
	if err != nil {
		return nil
	}
	if strings.Contains(response, "callback"){
		return nil
	}
	temp := strings.Split(response, "&")[0]
	accesstoken := strings.Split(temp, "=")[1]
	return map[string]interface{}{"access_token" : accesstoken}
}

func (oauth *QQOAuth) GetOpenid(accesstoken string) map[string]interface{}{
	request:= httplib.Get(qq_openid_url)
	request.Param("access_token",accesstoken)
	request.Param("unionid","1")
	responseStr, _ := request.String()
	var response  map[string]interface{}
	json.Unmarshal([]byte(responseStr[10:len(responseStr)-3]),&response)
	return response
}

func (oauth *QQOAuth) GetUserinfo(accesstoken string, openid string) map[string]interface{}{
	request:= httplib.Get(qq_getuserinfo_url)
	request.Param("access_token",accesstoken)
	request.Param("oauth_consumer_key",oauth.appkey)
	request.Param("openid",openid)
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func (oauth *QQOAuth) Authorize(code string) AuthorizeResult{
	accesstokenResponse := oauth.GetAccesstoken(code)
	if accesstokenResponse == nil{
		return AuthorizeResult{false, nil}
	}
	accesstoken := accesstokenResponse["access_token"].(string)
	openidResponse := oauth.GetOpenid(accesstoken)
	if _, ok := openidResponse["error"]; ok {  //获取openid接口返回错误
		return AuthorizeResult{false, nil}
	}
	openid := openidResponse["openid"].(string)
	unionid :=  openidResponse["unionid"].(string)

	getuserinfoResult := oauth.GetUserinfo(accesstoken, openid)
	if getuserinfoResult == nil{
		return AuthorizeResult{false, nil}
	}
	var sex int
	gender, ok := getuserinfoResult["gender"]
	if !ok {
		sex = 1
	}
	if gender.(string) == "女" {
		sex = 2
	}else{
		sex = 1
	}
	return AuthorizeResult{true, map[string]interface{}{
		"nickname":getuserinfoResult["nickname"].(string),
		"openid":openid,
		"sex":sex,
		"headimgurl":getuserinfoResult["figureurl_qq_1"].(string),   // QQ头像 40x40尺寸
		"unionid":unionid}}
}

func (oauth *QQOAuth) InitOAuth(){
	oauth.appkey = QQappkey
	oauth.appsecret = QQappsecret
	oauth.redirect_url = QQRedirectUrl
}

func init(){
	ReisterPlatform("qq", qqOAuth)
}
