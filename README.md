# simpleoauth

Golang实现几大主流平台的oauth2.0认证（目前仅支持QQ，微信，微博）

使用方法：
在config.go中配置好相关信息后

    package main

     import (
     	"fmt"
     	"simpleoauth"
     )

     func main() {
     	m, _ := simpleoauth.NewManager("qq")
     	result := m.Authorize("此处填入认证通过后，第三方平台跳转带回来的CODE")
     	fmt.Println(result)
     }
     
返回结果为：

    type AuthorizeResult struct{
        Result bool                         //认证成功或失败
        Userinfo map[string]interface{}     //返回用户信息等基本资料
    }