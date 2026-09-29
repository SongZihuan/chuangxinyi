package main

import (
	"fmt"
	user "gitee.com/wuntsong-auth/backend/src/service/v1"
)

func main() {
	fmt.Println("Start user auth backend...")
	user.CmdMain()
	fmt.Println("Bye~")
}
