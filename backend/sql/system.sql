CREATE TABLE IF NOT EXISTS access_record (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `request_id_prefix` VARCHAR(100) NOT NULL,  -- 请求ID的前缀
    `server_name` VARCHAR(100) NOT NULL, -- 服务器peerName FIRST
    `user_id` INT UNSIGNED,  -- 用户ID LAST
    `user_uid` VARCHAR(64),  -- 用户UID LAST
    `user_token` TEXT,  -- 用户TOKEN LAST
    `role_id` INT UNSIGNED,  -- 角色ID LAST
    `role_name` VARCHAR(20),  -- 角色名称 LAST
    `role_sign` VARCHAR(100),  -- 角色标志 LAST
    `web_id` INT UNSIGNED,  -- 站点名称 LAST
    `web_name` VARCHAR(100),  -- 站点标志 LAST
    `requests_web_id` INT UNSIGNED,  -- 站点名称 LAST
    `requests_web_name` VARCHAR(100),  -- 站点标志 LAST
    `ip` VARCHAR(100) NOT NULL,  -- X-Real-IP FIRST
    `geo_code` VARCHAR(10) NOT NULL,  -- 行政区划编码 FIRST
    `geo` VARCHAR(100) NOT NULL,  -- 地理位置 FIRST
    `scheme` VARCHAR(20) NOT NULL,  -- 协议 FIRST
    `method` VARCHAR(20) NOT NULL,  -- 调用方法 FIRST
    `host` VARCHAR(500) NOT NULL,  -- host FIRST
    `path` VARCHAR(1000) NOT NULL,  -- 调用path FIRST
    `query` JSON NOT NULL,  -- query FIRST
    `content_type` VARCHAR(100) NOT NULL,  -- Content-Type FIRST
    `requests_body` TEXT NOT NULL,  -- 请求的body FIRST
    `response_body` TEXT,  -- 返回的body LAST
    `response_body_error` VARCHAR(2000),  -- 返回的body写入错误 LAST
    `requests_header` JSON NOT NULL,  -- 请求的header FIRST
    `response_header` JSON,  -- 返回的header LAST
    `status_code` INT,  -- 状态码 LAST
    `panic_error` VARCHAR(2000),  -- 触发的宕机错误 LAST
    `message` JSON,  -- message LAST
    `use_time` INT,  -- 消耗时间（毫秒） LAST
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `start_at` DATETIME NULL DEFAULT NULL,  -- 开始时间 LAST
    `end_at` DATETIME NULL DEFAULT NULL,  -- 结束事件 LAST
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS token_record (
    `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
    `token_type` INT UNSIGNED NOT NULL,  -- UserToken LoginToken
    `token` TEXT NOT NULL,  -- 用户ID
    `type` INT UNSIGNED NOT NULL,  -- 创建、变更地点、删除
    `data` JSON NOT NULL,  -- 信息
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end


CREATE TABLE IF NOT EXISTS website (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `uid` VARCHAR(64) NOT NULL,
    `name` VARCHAR(100) NOT NULL,  -- 名字
    `pubkey` VARCHAR(1000) NOT NULL,  -- 公钥
    `describe` VARCHAR(1000) NOT NULL,  -- 介绍
    `keymap` VARCHAR(1000) NOT NULL,  -- 键值对
    `agreement` TEXT NOT NULL,  -- 协议
    `permission` VARCHAR(128) NOT NULL,
    `status` INT NOT NULL,  -- 状态
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE website MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS website_ip (
   `id` INT UNSIGNED AUTO_INCREMENT,
   `website_id` INT UNSIGNED NOT NULL,
   `ip` VARCHAR(100) NOT NULL,
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS website_domain (
  `id` INT UNSIGNED AUTO_INCREMENT,
  `website_id` INT UNSIGNED NOT NULL,
  `domain` VARCHAR(200) NOT NULL,
  `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
  `delete_at` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS agreement (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `aid`       VARCHAR(100) NOT NULL,
    `content` TEXT NOT NULL,
    `create_at` DATETIME     NOT NULL DEFAULT NOW(), -- 创建时间
    `update_at` DATETIME     NOT NULL DEFAULT NOW(), -- 更新时间
    `delete_at` DATETIME     NULL     DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE agreement MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL; -- end

CREATE TABLE IF NOT EXISTS footer (
     `id` INT UNSIGNED AUTO_INCREMENT,
     `copyright` VARCHAR(100) NOT NULL,
     `icp1` VARCHAR(100) NOT NULL,
     `icp2` VARCHAR(100) NOT NULL,
     `gongan` VARCHAR(100) NOT NULL,
     `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
     PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS oss_file (
      `id` INT UNSIGNED AUTO_INCREMENT,
      `fid` VARCHAR(100) NOT NULL,
      `key` VARCHAR(100) NOT NULL,
      `media_type` VARCHAR(100) NOT NULL,  -- 媒体类型
      `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
      `delete_at` DATETIME NULL DEFAULT NULL,
      PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS website_funding (
  `id` INT UNSIGNED AUTO_INCREMENT,  -- 记录ID
  `web_id` INT UNSIGNED NOT NULL,  -- 用户ID
  `type` INT NOT NULL,
  `funding_id` VARCHAR(64) NOT NULL,
  `profit` INT NOT NULL,  -- 获利
  `expenditure` INT NOT NULL,  -- 指出
  `year` INT NOT NULL,  -- 年份
  `month` INT NOT NULL,  -- 月份
  `day` INT NOT NULL,  -- 日期
  `remark` TEXT NOT NULL,
  `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
  `pay_at` DATETIME NOT NULL,
  `delete_at` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

CREATE TABLE IF NOT EXISTS application (
   `id` INT UNSIGNED AUTO_INCREMENT,
   `name` VARCHAR(200) NOT NULL,
   `describe` VARCHAR(200) NOT NULL,
   `web_id` INT UNSIGNED NOT NULL,
   `url` VARCHAR(500) NOT NULL,
   `icon` VARCHAR(100) NOT NULL,
   `status` INT NOT NULL,  -- 状态
   `sort` INT UNSIGNED NOT NULL,
   `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
   `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
   `delete_at` DATETIME NULL DEFAULT NULL,
   PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE application MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS policy (
    `id` INT UNSIGNED UNIQUE AUTO_INCREMENT,  -- 权限ID
    `name` VARCHAR(20) NOT NULL,  -- 权限名称
    `sign` VARCHAR(100) NOT NULL,  -- 权限表示
    `describe` VARCHAR(200) NOT NULL,  -- 权限描述
    `sort` INT UNSIGNED NOT NULL,  -- 排序
    `is_anonymous` BOOL NOT NULL,  -- 是否匿名权限
    `is_user` BOOL NOT NULL,  -- 是否一般用户权限
    `status` INT NOT NULL,  -- 状态：停用、启用
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE policy MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS url_path (
    `id` INT UNSIGNED UNIQUE AUTO_INCREMENT,  -- 权限ID
    `describe` VARCHAR(200) NOT NULL,  -- 路由描述
    `path` VARCHAR(200) NOT NULL,  -- 路由
    `mode` INT NOT NULL,  -- 前缀匹配、正则匹配、全匹配
    `authentication` BOOL NOT NULL,  -- 是否启用鉴权
    `double_check` BOOL NOT NULL,
    `cors_mode` INT NOT NULL,  -- all, not website, website allow
    `admin_mode` INT NOT NULL, -- not admin, website admin, normal admin
    `busy_mode` INT NOT NULL,  -- busy模式
    `busy_count` INT NOT NULL,  -- busy计数
    `captcha_mode` INT NOT NULL,
    `status` INT NOT NULL,  -- 状态：停用、启用
    `is_or_policy` BOOL NOT NULL,  -- 权限是or关系还是and关系 or-true and-false
    `permission` VARCHAR(128) NOT NULL,
    `sub_policy` BIGINT NOT NULL,
    `method` BIGINT NOT NULL,
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE url_path MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS website_policy (
  `id` INT UNSIGNED UNIQUE AUTO_INCREMENT,  -- 权限ID
  `name` VARCHAR(20) NOT NULL,  -- 权限名称
  `sign` VARCHAR(100) NOT NULL,  -- 权限表示
  `describe` VARCHAR(200) NOT NULL,  -- 权限描述
  `sort` INT UNSIGNED NOT NULL,  -- 排序
  `status` INT NOT NULL,  -- 状态：停用、启用
  `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
  `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
  `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
  PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE website_policy MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS website_url_path (
    `id` INT UNSIGNED UNIQUE AUTO_INCREMENT,  -- 权限ID
    `describe` VARCHAR(200) NOT NULL,  -- 路由描述
    `path` VARCHAR(200) NOT NULL,  -- 路由
    `mode` INT NOT NULL,  -- 前缀匹配、正则匹配、全匹配
    `status` INT NOT NULL,  -- 状态：停用、启用
    `is_or_policy` BOOL NOT NULL,  -- 权限是or关系还是and关系 or-true and-false
    `permission` VARCHAR(128) NOT NULL,
    `method` BIGINT NOT NULL,
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE website_url_path MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS role (
    `id` INT UNSIGNED UNIQUE AUTO_INCREMENT,  -- 角色ID
    `name` VARCHAR(20) NOT NULL,  -- 角色名称
    `sign` VARCHAR(100) NOT NULL,  -- role标识符
    `describe` VARCHAR(200) NOT NULL,  -- 角色描述
    `belong` INT UNSIGNED,  -- 归属的站点
    `status` INT NOT NULL,  -- 状态：停用、启用
    `permissions` VARCHAR(128) NOT NULL,  -- 权限
    `not_delete` BOOL NOT NULL,  -- 不可删除
    `not_change_sign` BOOL NOT NULL,  -- 不可修改sign
    `not_change_permissions` BOOL NOT NULL,  -- 不可修改permissions和belong
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE role MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS menu (
    `id` INT UNSIGNED UNIQUE AUTO_INCREMENT,  -- 菜单ID
    `sort` INT UNSIGNED NOT NULL,
    `describe` VARCHAR(200) NOT NULL,
    `father_id` INT UNSIGNED,  -- 菜单父ID
    `name` VARCHAR(100) NOT NULL,  -- 菜单名字
    `path` VARCHAR(100) NOT NULL,  -- 菜单路径
    `title` VARCHAR(100) NOT NULL,  -- 菜单标题
    `icon` VARCHAR(100) NOT NULL,  -- 菜单图标
    `redirect` VARCHAR(100),  -- 菜单跳转
    `superior` VARCHAR(100) NOT NULL,
    `category` INT NOT NULL,
    `component` VARCHAR(100) NOT NULL,  -- 组件
    `component_alias` VARCHAR(100) NOT NULL,  -- 组件
    `meta_link` VARCHAR(100),  -- 菜单链接
    `type` INT NOT NULL,
    `is_link` BOOL NOT NULL,
    `is_hide` BOOL NOT NULL,
    `is_keepalive` BOOL NOT NULL,
    `is_affix` BOOL NOT NULL,
    `is_iframe`  BOOL NOT NULL,
    `btn_power` VARCHAR(100) NOT NULL,
    `is_or_policy` BOOL NOT NULL,  -- 权限是or关系还是and关系 or-true and-false
    `status` INT NOT NULL,
    `policy` VARCHAR(128) NOT NULL,
    `sub_policy` BIGINT NOT NULL,
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `update_at` DATETIME NOT NULL DEFAULT NOW(),  -- 更新时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end

ALTER TABLE menu MODIFY update_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() NOT NULL;  -- end

CREATE TABLE IF NOT EXISTS announcement (
    `id` INT UNSIGNED UNIQUE AUTO_INCREMENT,  -- 角色菜单ID
    `sort` INT NOT NULL,  -- 排序
    `title` VARCHAR(100) NOT NULL,  -- 标题
    `content` VARCHAR(1000) NOT NULL,  -- 内容
    `create_at` DATETIME NOT NULL DEFAULT NOW(),  -- 创建时间
    `start_at` DATETIME NOT NULL,  -- 开始时间
    `stop_at` DATETIME NOT NULL,  -- 结束时间
    `delete_at` DATETIME NULL DEFAULT NULL,  -- 删除时间
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4  COLLATE utf8mb4_bin; -- end