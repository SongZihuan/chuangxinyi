package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	errors "github.com/wuntsong-org/wterrors"
	"strings"
	"time"
)

func ClearDataByDeleteAt(ctx context.Context, tableName string, deleteTime time.Time) (sql.Result, errors.WTError) {
	query := fmt.Sprintf("delete from %s where `delete_at` < ?", tableName)
	res, err := MySQLConn.ExecCtx(ctx, query, deleteTime)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return res, nil
}

func ClearDataByCreateAt(ctx context.Context, tableName string, deleteTime time.Time) (sql.Result, errors.WTError) {
	query := fmt.Sprintf("delete from %s where `create_at` < ?", tableName)
	res, err := MySQLConn.ExecCtx(ctx, query, deleteTime)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return res, nil
}

func ClearData() errors.WTError {
	tableCount := len(config.BackendConfig.MySQL.ClearDeleteAt) + len(config.BackendConfig.MySQL.ClearCreateAt)

	logger.Logger.WXInfo("开始清理数据库，计划清理 %d 张表", tableCount)

	table := make(map[string]bool, tableCount)
	errorTables := 0
	processTables := 0
	processRows := int64(0)

	for _, t := range config.BackendConfig.MySQL.ClearDeleteAt {
		var err error

		tableName := strings.TrimSpace(t.TableName)

		if len(tableName) == 0 {
			continue
		}

		exists, ok := table[tableName]
		if ok && exists {
			continue
		}

		table[tableName] = true

		saveDay := t.SaveDay

		if saveDay <= 0 {
			saveDay = 180
		}

		ctx := context.Background()
		res, err := ClearDataByDeleteAt(ctx, tableName, time.Now().Add(-time.Hour*24*time.Duration(saveDay)))
		if err != nil {
			errorTables += 1
			logger.Logger.WXInfo("清理数据库表（%s）时出错：%s", tableName, err.Error())
			continue
		}

		count, err := res.RowsAffected()
		if err != nil {
			errorTables += 1
			logger.Logger.WXInfo("清理数据库表（%s）时出错：%s", tableName, err.Error())
			continue
		}

		if count > 0 {
			processTables += 1
			processRows += count
			logger.Logger.WXInfo("清理数据库表（%s）成功，清理了%d条数据", tableName, count)
		}
	}

	for _, t := range config.BackendConfig.MySQL.ClearCreateAt {
		var err error

		tableName := strings.TrimSpace(t.TableName)

		if len(tableName) == 0 {
			continue
		}

		exists, ok := table[tableName]
		if ok && exists {
			continue
		}

		table[tableName] = true
		saveDay := t.SaveDay

		if saveDay <= 0 {
			saveDay = 180
		}

		ctx := context.Background()
		res, err := ClearDataByCreateAt(ctx, tableName, time.Now().Add(-time.Hour*24*time.Duration(saveDay)))
		if err != nil {
			errorTables += 1
			logger.Logger.WXInfo("清理数据库表（%s）时出错：%s", tableName, err.Error())
			continue
		}

		count, err := res.RowsAffected()
		if err != nil {
			errorTables += 1
			logger.Logger.WXInfo("清理数据库表（%s）时出错：%s", tableName, err.Error())
			continue
		}

		if count > 0 {
			processTables += 1
			processRows += count
			logger.Logger.WXInfo("清理数据库表（%s）成功，清理了%d条数据", tableName, count)
		}
	}

	logger.Logger.WXInfo("清理数据库结束，清理出错 %d 张表，清理成功 %d 张表，清理成功删除 %d 行", errorTables, processTables, processRows)

	return nil
}
