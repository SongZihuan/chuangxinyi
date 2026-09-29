package main

import (
	"flag"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	initall "gitee.com/wuntsong-auth/backend/src/init"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/utils"
)

var configFile = flag.String("f", "etc", "the config path")

func main() {
	fmt.Println("Start sql clear...")
	CmdMain()
	fmt.Println("Bye~")
}

func CmdMain() {
	flag.Parse()

	err := config.InitBackendConfigViper(*configFile, "AUTH_")
	utils.MustNotError(err)

	err = initall.InitSqlClear()
	utils.MustNotError(err)

	err = mysql.ClearData()
	utils.MustNotError(err)
}
