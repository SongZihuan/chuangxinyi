CREATE TABLE IF NOT EXISTS user (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 用户ID
    `uid` VARCHAR(64) NOT NULL,  -- 用户展示ID
    `status` INT NOT NULL,  -- WAIT_REG(等待注册流程), NORMAL(正常使用), BANNED(封禁)
    `signin` BOOL NOT NULL DEFAULT FALSE,  -- 单点登录
    `son_level` INT UNSIGNED NOT NULL,  -- 儿子层级
    `father_id` INT UNSIGNED,  -- 父亲ID
    `root_father_id` INT UNSIGNED,  -- 根父亲ID
    `invite_id` INT UNSIGNED,  -- 邀请ID
    `wallet_id` INT UNSIGNED NOT NULL,  -- 钱包ID
    `token_expiration` INT UNSIGNED NOT NULL,  -- token持续时间
    `role_id` INT NOT NULL,  -- 角色ID
    `is_admin` BOOL NOT NULL DEFAULT FALSE,  -- 是否根admin
    `remark` TEXT NOT NULL,
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE user MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS phone (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `phone` VARCHAR(20) NOT NULL,  -- 用户邮箱
     `is_delete` BOOL NOT NULL,  -- 是否删除
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS wallet (
      `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
      `balance` INT UNSIGNED NOT NULL,  -- 余额
      `wait_balance` INT UNSIGNED NOT NULL,  -- 待入账余额
      `cny` INT UNSIGNED NOT NULL,  -- 实际充值金额（在平台实际充值的钱）
      `not_billed` INT NOT NULL,  -- 未开票金额（负数表示欠票）
      `billed` INT UNSIGNED NOT NULL,  -- 可开票总金额
      `has_billed` INT UNSIGNED NOT NULL,  -- 已开票总金额
      `withdraw` INT UNSIGNED NOT NULL,  -- 总共可提现
      `wait_withdraw` INT UNSIGNED NOT NULL,  -- 待入账可提现
      `not_withdraw` INT UNSIGNED NOT NULL,  -- 未提现
      `has_withdraw` INT UNSIGNED NOT NULL,  -- 已提现
      `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
      `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
      `delete_at` DATETIME NULL DEFAULT NULL,
      PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE wallet MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS nickname (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `nickname` VARCHAR(50),  -- 用户昵称 NULL为删除
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS header (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `header` VARCHAR(64),  -- 用户头像 NULL为删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS email (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `email` VARCHAR(50),  -- 用户邮箱 NULL为删除
    `is_delete` BOOL NOT NULL,  -- 是否删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS wechat (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `open_id` VARCHAR(128),  -- OpenID
    `union_id` VARCHAR(128),  -- UnionID
    `fuwuhao` VARCHAR(128),  -- 服务号OpenID
    `nickname` VARCHAR(100),  -- 用户微信昵称
    `headimgurl` VARCHAR(200),  -- 头像
    `is_delete` BOOL NOT NULL,  -- 是否删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS wxrobot (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `webhook` VARCHAR(200),  -- 微信机器人webhook
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS password (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `passwordHash` VARCHAR(70),  -- 密码哈希（Hash256，两次） NULL为删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS username (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `username` VARCHAR(200),  -- 账号名称
    `is_delete` BOOL NOT NULL,  -- 是否删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS secondfa (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `secret` VARCHAR(150),  -- 用户2FA密钥 NULL为删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS title (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `name` VARCHAR(30),  -- 姓名, 公司名
     `tax_id` VARCHAR(20),  -- 税号, 身份证号
     `bank_id` VARCHAR(20),  -- 银行卡号
     `bank` VARCHAR(30),  -- 开户行
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS address (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `name` VARCHAR(30),  -- 收件人
     `phone` VARCHAR(20),  -- 收件人手机
     `email` VARCHAR(50),  -- 收件人邮箱
     `country` VARCHAR(20),  -- 国家
     `province` VARCHAR(20),  -- 省份
     `city` VARCHAR(20),  -- 城市
     `district` VARCHAR(20),  -- 区县
     `country_code` VARCHAR(10),
     `province_code` VARCHAR(10),  -- 省份代码
     `city_code` VARCHAR(10),  -- 城市代码
     `district_code` VARCHAR(10),  -- 区县代码
     `address` VARCHAR(100),  -- 详细地址
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS idcard (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `user_name` VARCHAR(20) NOT NULL,  -- 用户姓名
    `user_id_card` VARCHAR(20) NOT NULL,  -- 用户身份证号码
    `phone` VARCHAR(20),
    `idcard_key` VARCHAR(100),  -- 原件地址
    `idcard_back_key` VARCHAR(100),  -- 原件地址
    `face_check_id` VARCHAR(64),
    `is_company` BOOL NOT NULL,  -- 是否企业
    `is_delete` BOOL NOT NULL,  -- 是否删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE idcard MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS company (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `legal_person_name` VARCHAR(20)NOT NULL,  -- 用户法人姓名
    `legal_person_id_card` VARCHAR(20)NOT NULL,  -- 用户法人身份证号码
    `company_name` VARCHAR(30)NOT NULL,  -- 用户公司名
    `company_id` VARCHAR(20)NOT NULL,  -- 用户公司统一社会信用代码
    `license_key` VARCHAR(100),  -- 原件地址
    `idcard_key` VARCHAR(100),  -- 原件地址
    `idcard_back_key` VARCHAR(100),  -- 原件地址
    `face_check_id` VARCHAR(64),
    `is_delete` BOOL NOT NULL,  -- 是否删除
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE company MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS homepage (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `introduction` TEXT,  -- 用户简介
    `address` VARCHAR(200),  -- 联系地址
    `phone` VARCHAR(20),  -- 联系电话
    `email` VARCHAR(100),  -- 联系邮件
    `wechat` VARCHAR(50),  -- 微信号
    `qq` VARCHAR(50),  -- qq号
    `man` BOOL,  -- 是否男性
    `link` VARCHAR(200),  -- 外部连接
    `company` VARCHAR(50),  -- 行业
    `industry` VARCHAR(50),  -- 行业
    `position` VARCHAR(50),  -- 职位
    `close` BOOL NOT NULL,
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS pay (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `wallet_id` INT UNSIGNED NOT NULL,  -- 钱包ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `subject` VARCHAR(200) NOT NULL,  -- 产品名字
     `pay_way` VARCHAR(100) NOT NULL,  -- 支付方式
     `pay_id` VARCHAR(64) NOT NULL,
     `cny` INT NOT NULL,  -- 支付金额，单位：分
     `get` INT NOT NULL,  -- 获得额度
     `coupons_id` INT UNSIGNED,  -- 优惠券
     `trade_no` VARCHAR(100),  -- 支付宝或微信交易流水号
     `buyer_id` VARCHAR(100),  -- 购买者ID
     `trade_status` INT UNSIGNED NOT NULL,  -- WAIT, SUCCESS, FINISH, CLOSE, REFUND
     `balance` INT UNSIGNED,  -- 充值后的余额
     `remark` TEXT NOT NULL,
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `pay_at` DATETIME NULL DEFAULT NULL,  -- 支付时间
     `refund_at` DATETIME NULL DEFAULT NULL,  -- 退款时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS back (
   `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
   `wallet_id` INT UNSIGNED NOT NULL,  -- 钱包ID
   `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
   `subject` VARCHAR(200) NOT NULL,  -- 产品名字
   `back_id` VARCHAR(64) NOT NULL,
   `get` INT NOT NULL,  -- 获得额度
   `balance` INT UNSIGNED NOT NULL,  -- 充值后的余额
   `can_withdraw` BOOL NOT NULL,  -- 可提现
   `supplier_id` INT UNSIGNED NOT NULL,  -- 返现方ID（website ID）
   `supplier` VARCHAR(200) NOT NULL,
   `remark` TEXT NOT NULL,
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS withdraw (
      `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
      `wallet_id` INT UNSIGNED NOT NULL,  -- 钱包ID
      `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
      `withdraw_id` VARCHAR(64) NOT NULL,
      `withdraw_way` VARCHAR(100) NOT NULL,
      `name` VARCHAR(20) NOT NULL,  -- 提现人姓名
      `alipay_login_id` VARCHAR(50),  -- 支付宝快捷提现使用
      `wechatpay_open_id` VARCHAR(128),  -- 微信快捷提现使用
      `wechatpay_union_id` VARCHAR(128),  -- 微信快捷提现使用
      `wechatpay_nickname` VARCHAR(100),  -- 微信快捷提现使用
      `cny` INT NOT NULL,
      `balance` INT UNSIGNED,  -- 支付后的余额
      `order_id` VARCHAR(64),  -- 支付宝和微信用
      `pay_fund_order_id` VARCHAR(64),
      `remark` TEXT NOT NULL,
      `status` INT NOT NULL,  -- 等待提现，已提现，已取消，
      `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
      `withdraw_at` DATETIME NOT NULL ,  -- 支付时间
      `pay_at` DATETIME NULL DEFAULT NULL,  -- 支付时间
      `delete_at` DATETIME NULL DEFAULT NULL,
      PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS defray (
       `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
       `wallet_id` INT UNSIGNED,  -- 钱包ID
       `user_id` INT UNSIGNED,  -- 用户ID
       `owner_id` INT UNSIGNED,  -- 获得者ID
       `defray_id` VARCHAR(64) NOT NULL,
       `subject` VARCHAR(200) NOT NULL,  -- 产品名字
       `price` INT NOT NULL,  -- 支付金额
       `coupons_id` INT UNSIGNED,
       `unit_price` INT NOT NULL,
       `quantity` INT NOT NULL,
       `describe` VARCHAR(500) NOT NULL,
       `supplier_id` INT UNSIGNED NOT NULL,  -- 制造方ID（website ID）
       `supplier` VARCHAR(200) NOT NULL,
       `return_url` VARCHAR(200) NOT NULL,
       `real_price` INT UNSIGNED,  -- 实际支付金额
       `balance` INT UNSIGNED,  -- 支付后的余额
       `invite_pre` INT NOT NULL,
       `distribution_level_1` INT NOT NULL,  -- 1级分销
       `distribution_level_2` INT NOT NULL,  -- 2级分销
       `distribution_level_3` INT NOT NULL,  -- 3级分销
       `has_distribution` BOOL NOT NULL,  -- 是否一级分销
       `can_withdraw` BOOL NOT NULL,  -- 能否提现
       `remark` TEXT NOT NULL,
       `must_self_defray` BOOL NOT NULL,  -- 必须自己支付
       `return_reason` VARCHAR(200),  -- 退款原因
       `status` INT NOT NULL,  -- 等待支付, 已支付, 已退款
       `return_day_limit` INT NOT NULL,  -- 可退款日期限制
       `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
       `defray_at` DATETIME NULL DEFAULT NULL,  -- 支付时间
       `last_return_at` DATETIME NULL DEFAULT NULL,  -- 最后可退款时间
       `return_at` DATETIME NULL DEFAULT NULL,  -- 支付时间
       `delete_at` DATETIME NULL DEFAULT NULL,
       PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS invoice (
   `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
   `wallet_id` INT UNSIGNED NOT NULL,  -- 钱包ID
   `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
   `invoice_id` VARCHAR(64) NOT NULL,  -- 发票id
   `amount` INT UNSIGNED NOT NULL,  -- 金额
   `type` INT NOT NULL,  -- 个人普票，企业普票，企业专票

   `name` VARCHAR(30),  -- 姓名, 公司名
   `tax_id` VARCHAR(20),  -- 税号, 身份证号
   `bank_id` VARCHAR(20),  -- 银行卡号
   `bank` VARCHAR(30),  -- 开户行

   `recipient` VARCHAR(30),  -- 收件人
   `phone` VARCHAR(20),  -- 收件人手机
   `email` VARCHAR(50),  -- 收件人邮箱
   `country` VARCHAR(20),  -- 国家
   `province` VARCHAR(20),  -- 省份
   `city` VARCHAR(20),  -- 城市
   `district` VARCHAR(20),  -- 区县
   `address` VARCHAR(100),  -- 详细地址

   `invoice_number` VARCHAR(20),  -- 发票号码
   `invoice_code` VARCHAR(20),  -- 发票代码
   `invoice_check_code` VARCHAR(30),  -- 校验码
   `issuer_at` DATETIME NULL DEFAULT NULL,  -- 开票年
   `invoice_key` VARCHAR(100),  -- 发票key

   `red_invoice_number` VARCHAR(20),  -- 发票号码
   `red_invoice_code` VARCHAR(20),  -- 发票代码
   `red_invoice_check_code` VARCHAR(30),  -- 校验码
   `red_issuer_at` DATETIME NULL DEFAULT NULL,  -- 开票年
   `red_invoice_key` VARCHAR(100),  -- 发票keyr

   `remark` TEXT NOT NULL,
   `status` INT NOT NULL,  -- 待开票、已开票、已退票、错误、红冲
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `billing_at` DATETIME NULL DEFAULT NULL,  -- 开票时间
   `return_at` DATETIME NULL DEFAULT NULL,  -- 退票时间
   `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS wallet_record (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `wallet_id` INT UNSIGNED NOT NULL,  -- 钱包ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `type` INT NOT NULL,
    `funding_id` VARCHAR(64) NOT NULL,
    `reason` VARCHAR(200) NOT NULL,

    `balance` INT UNSIGNED NOT NULL,  -- 余额
    `wait_balance` INT UNSIGNED NOT NULL,  -- 待入账余额
    `cny` INT UNSIGNED NOT NULL,  -- 余额
    `not_billed` INT UNSIGNED NOT NULL,  -- 未开票金额 （可能为负数，表示欠票）
    `billed` INT UNSIGNED NOT NULL,  -- 可开票总金额
    `has_billed` INT UNSIGNED NOT NULL,  -- 已开票总金额
    `withdraw` INT UNSIGNED NOT NULL,  -- 总共可提现
    `wait_withdraw` INT UNSIGNED NOT NULL,  -- 待入账余额
    `not_withdraw` INT UNSIGNED NOT NULL,  -- 未提现
    `has_withdraw` INT UNSIGNED NOT NULL,  -- 已提现

    `before_balance` INT UNSIGNED NOT NULL,  -- 余额
    `before_wait_balance` INT UNSIGNED NOT NULL,  -- 待入账余额
    `before_cny` INT UNSIGNED NOT NULL,  -- 余额
    `before_not_billed` INT UNSIGNED NOT NULL,  -- 未开票金额 （可能为负数，表示欠票）
    `before_billed` INT UNSIGNED NOT NULL,  -- 可开票总金额
    `before_has_billed` INT UNSIGNED NOT NULL,  -- 已开票总金额
    `before_withdraw` INT UNSIGNED NOT NULL,  -- 总共可提现
    `before_wait_withdraw` INT UNSIGNED NOT NULL,  -- 待入账余额
    `before_not_withdraw` INT UNSIGNED NOT NULL,  -- 未提现
    `before_has_withdraw` INT UNSIGNED NOT NULL,  -- 已提现

    `remark` TEXT NOT NULL,
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 支付时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS message (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
    `title` VARCHAR(100) NOT NULL,
    `content` TEXT NOT NULL,
    `sender` VARCHAR(100) NOT NULL,  -- 发送方
    `sender_id` INT UNSIGNED NOT NULL,  -- 发送方ID（website ID）
    `sender_link` VARCHAR(200),  -- 发送方链接
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `read_at` DATETIME NULL DEFAULT NULL,  -- 阅读时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS sms_message (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `phone` VARCHAR(20) NOT NULL,  -- 用户ID
    `sig` VARCHAR(100) NOT NULL,
    `template` VARCHAR(100) NOT NULL,
    `template_param` JSON NOT NULL,
    `sender_id` INT UNSIGNED NOT NULL,  -- 发送方ID（website ID）
    `success` bool NOT NULL,
    `error_msg` VARCHAR(200),
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS email_message (
   `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
   `email` VARCHAR(50) NOT NULL,  -- 用户ID
   `subject` VARCHAR(100) NOT NULL,
   `content` TEXT NOT NULL,
   `sender` VARCHAR(100) NOT NULL,
   `sender_id` INT UNSIGNED NOT NULL,  -- 发送方ID（website ID）
   `success` bool NOT NULL,
   `error_msg` VARCHAR(200),
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS fuwuhao_message (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `open_id` VARCHAR(128) NOT NULL,  -- 用户ID
     `template` VARCHAR(100) NOT NULL,
     `url` VARCHAR(200) NOT NULL,
     `val` JSON NOT NULL,
     `sender_id` INT UNSIGNED NOT NULL,  -- 发送方ID（website ID）
     `success` bool NOT NULL,
     `error_msg` VARCHAR(200),
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS wxrobot_message (
   `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
   `webhook` VARCHAR(200) NOT NULL,  -- 用户ID
   `text` VARCHAR(1000) NOT NULL,
   `at_all` bool NOT NULL,
   `sender_id` INT UNSIGNED NOT NULL,  -- 发送方ID（website ID）
   `success` bool NOT NULL,
   `error_msg` VARCHAR(200),
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS audit (
   `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
   `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
   `content` TEXT NOT NULL,  -- 审计内容
   `from` VARCHAR(100) NOT NULL,  -- 条目来源
   `from_id` INT UNSIGNED NOT NULL,  -- 制造方ID（website ID）
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS uncle (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `uncle_id` INT UNSIGNED NOT NULL,  -- 叔叔ID
     `uncle_tag` VARCHAR(200) NOT NULL,  -- 叔叔标志（用于在未获取叔叔授权之前标识叔叔）
     `status` INT NOT NULL,
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id),
     UNIQUE KEY user_uncle (user_id, uncle_id, delete_at)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS work_order (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `uid` VARCHAR(64) NOT NULL,  -- 工单ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `title` VARCHAR(100) NOT NULL,  -- 工单ID
     `from` VARCHAR(100) NOT NULL,  -- 工单ID
     `from_id` INT UNSIGNED NOT NULL,  -- 制造方ID（website ID）
     `remark` TEXT NOT NULL,
     `status` INT NOT NULL,  -- (等待用户回复，等待网站回复，已完成)
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `last_reply_at`  DATETIME NULL DEFAULT NULL,  -- 上次回复时间
     `finish_at`  DATETIME NULL DEFAULT NULL,  -- 完成时间时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS work_order_communicate (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `order_id` INT UNSIGNED NOT NULL,  -- 工单数字ID
    `content` TEXT NOT NULL,  -- 沟通内容
    `from` INT NOT NULL,  -- (来自用户，来自网站)
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS work_order_communicate_file (
      `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
      `order_id` INT UNSIGNED NOT NULL,  -- 工单数字ID
      `communicate_id` INT UNSIGNED NOT NULL,  -- 沟通记录ID
      `key` VARCHAR(100) NOT NULL,  -- 文件名称
      `fid` VARCHAR(64) NOT NULL,  -- 文件ID
      `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
      `delete_at` DATETIME NULL DEFAULT NULL,
      PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS discount (
     `id` INT UNSIGNED AUTO_INCREMENT,
     `name` VARCHAR(100) NOT NULL,
     `describe` TEXT NOT NULL,
     `short_describe` VARCHAR(200) NOT NULL,
     `type` INT UNSIGNED NOT NULL,  -- 送优惠券，送额度
     `quota` JSON NOT NULL,  -- 内容物表格
     `day_limit` INT UNSIGNED,  -- 日限购
     `month_limit` INT UNSIGNED,  -- 月限购
     `year_limit` INT UNSIGNED,  -- 年限购
     `limit` INT UNSIGNED,  -- 限购
     `need_verify` BOOL NOT NULL,
     `need_company` BOOL NOT NULL,
     `need_user_origin` BOOL NOT NULL,
     `need_company_origin` BOOL NOT NULL,
     `need_user_face`    BOOL NOT NULL,
     `need_company_face` BOOL NOT NULL,
     `show` BOOL NOT NULL,  -- 是否显示
     `remark` TEXT NOT NULL,
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
     `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
     PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE discount MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS discount_buy (
   `id` INT UNSIGNED AUTO_INCREMENT,
   `user_id` INT UNSIGNED NOT NULL,
   `discount_id` INT UNSIGNED NOT NULL,
   `name` VARCHAR(100) NOT NULL,
   `short_describe` VARCHAR(200) NOT NULL,
   `days` INT UNSIGNED NOT NULL,  -- 日购买量
   `month` INT UNSIGNED NOT NULL,  -- 日购买量
   `year` INT UNSIGNED NOT NULL,  -- 日购买量
   `all` INT UNSIGNED NOT NULL,  -- 日购买量
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
   PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS coupons (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `user_id` INT UNSIGNED NOT NULL,
    `type` INT UNSIGNED NOT NULL,  -- 充值送，满减，满打折
    `name` VARCHAR(100) NOT NULL,
    `content` JSON NOT NULL,  -- 内容物
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS oauth2_record (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `web_id` INT UNSIGNED NOT NULL,  -- 站点ID
     `web_name` VARCHAR(100) NOT NULL,  -- 站点名称
     `ip` VARCHAR(30) NOT NULL,  -- 登录IP
     `geo` VARCHAR(100) NOT NULL,  -- 登录地点
     `geo_code` VARCHAR(10) NOT NULL,  -- 登录地点代码
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `login_time` DATETIME NOT NULL,  -- 登录时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS oauth2_baned (
     `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
     `user_id` INT UNSIGNED NOT NULL,  -- 用户ID
     `web_id` INT UNSIGNED NOT NULL,  -- 站点ID
     `allow_login` BOOL NOT NULL,  -- 允许登录
     `allow_defray` BOOL NOT NULL,  -- 允许支付
     `allow_msg` BOOL NOT NULL,  -- 允许通信
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     `delete_at` DATETIME NULL DEFAULT NULL,
     PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS login_controller (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `user_id` INT UNSIGNED NOT NULL,
    `allow_phone` BOOL NOT NULL,
    `allow_email` BOOL NOT NULL,
    `allow_wechat` BOOL NOT NULL,
    `allow_password` BOOL NOT NULL,
    `allow_2fa` BOOL NOT NULL,
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS face_check (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `check_id` VARCHAR(64) NOT NULL,
    `certify_id` VARCHAR(100) NOT NULL,
    `name` VARCHAR(20) NOT NULL,
    `idcard` VARCHAR(20) NOT NULL,
    `status` INT NOT NULL,  -- 等待 通过 不通过
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end
