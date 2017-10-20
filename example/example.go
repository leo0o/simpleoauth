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
