package mysql

import (
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
)

var MySQL *sql.DB
var MySQLConn sqlx.SqlConn

func InitMysql() errors.WTError {
	if len(config.BackendConfig.MySQL.DSN) == 0 {
		return errors.Errorf("dsn must be given")
	}

	MySQL, err := sql.Open("mysql", config.BackendConfig.MySQL.DSN)
	if err != nil {
		fmt.Println("HHHH")
		return errors.WarpQuick(err)
	}

	err = MySQL.Ping()
	if err != nil {
		_ = MySQL.Close()
		MySQL = nil
		return errors.WarpQuick(err)
	}

	MySQLConn = sqlx.NewSqlConnFromDB(MySQL)

	err = ExecSqlFile()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func CloseMysql() {
	if MySQL == nil {
		return
	}

	_ = MySQL.Close()
	MySQLConn = nil
	MySQL = nil
}
