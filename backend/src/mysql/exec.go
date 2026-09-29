package mysql

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"os"
	"path"
	"strings"
)

const UserSQLFile = "user.sql"
const SystemSQLFile = "system.sql"

func ExecSqlFile() errors.WTError {
	err := execSqlFile(UserSQLFile)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = execSqlFile(SystemSQLFile)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func execSqlFile(file string) errors.WTError {
	sql, err := readSQLFile(path.Join(config.BackendConfig.MySQL.SQLFilePath, file))
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, s := range sql {
		_, err := MySQLConn.Exec(s)
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return nil
}

func readSQLFile(filepath string) ([]string, errors.WTError) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	defer utils.Close(f)

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	sql := strings.Split(string(data), "-- end")
	res := make([]string, 0, len(sql))
	for _, s := range sql {
		s = strings.Trim(s, "\r\n ")
		if len(s) != 0 {
			res = append(res, s)
		}
	}

	return res, nil
}
