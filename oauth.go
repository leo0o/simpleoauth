package simpleoauth

import (
	"fmt"
)

var oauthes= make(map[string]OAuth)

type OAuth interface{
	GetAccesstoken(code string) map[string]interface{}
	GetUserinfo(accesstoken string, openid string) map[string]interface{}
	Authorize(code string) AuthorizeResult
	InitOAuth()
}

func ReisterPlatform(name string, oauth OAuth){
	if  oauth == nil {
		panic("Register simpleoauth instance is nil")
	}
	_, dup := oauthes[name]
	if  dup{
		panic("The platform has registered already")
	}
	oauthes[name] = oauth
}

type Manager struct {
	oauth OAuth
}

func NewManager(platformName string)(*Manager, error){
	oauth, ok := oauthes[platformName]
	if !ok{
		return nil, fmt.Errorf("unknown platform %q", platformName)
	}
	oauth.InitOAuth()
	return &Manager{oauth}, nil
}

func (m *Manager)Authorize(code string) AuthorizeResult{
	return m.oauth.Authorize(code)
}

type AuthorizeResult struct{
	Result bool
	Userinfo map[string]interface{}
}
