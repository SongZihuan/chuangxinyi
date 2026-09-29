package db

import (
	"database/sql"
)

type OneInt struct {
	Res int64 `db:"res"`
}

type OneIntOrNull struct {
	Res sql.NullInt64 `db:"res"`
}

const (
	TimeCreateAt    = 1  // 创建时间（全部都有）
	TimePayAt       = 2  // 支付时间（支付记录）
	TimeDefrayAt    = 3  // 付款时间（消费记录）
	TimeBilledAt    = 4  // 开票时间（发票）
	TimeReturnAt    = 5  // 退票时间，退款时间（发票，消费记录）
	TimeReadAt      = 6  // 读取时间（站内信）
	TimeFinishAt    = 7  // 完成时间（工单）
	TimeLastReplyAt = 8  // 上次回复时间（工单）
	TimeStartAt     = 9  // 开始时间（访问记录）
	TimeEndAt       = 10 // 结束时间（访问记录）
	TimeLoginAt     = 11 // 登录时间（Oauth2记录）
	TimeWithdrawAt  = 12 // 提现时间（提现）
)

var timeMap = map[int64]string{
	TimeCreateAt:    "create_at",
	TimePayAt:       "pay_at",
	TimeDefrayAt:    "defray_at",
	TimeBilledAt:    "billed_at",
	TimeReturnAt:    "return_at",
	TimeReadAt:      "read_at",
	TimeFinishAt:    "finish_at",
	TimeLastReplyAt: "last_reply_at",
	TimeStartAt:     "start_at",
	TimeEndAt:       "end_at",
	TimeLoginAt:     "login_at",
	TimeWithdrawAt:  "withdraw_at",
}

func IsTimeType(t int64) bool {
	_, ok := timeMap[t]
	return ok
}
