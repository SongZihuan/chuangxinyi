-- 索引文件

ALTER TABLE user ADD CONSTRAINT idx_user_uid UNIQUE (uid);
ALTER TABLE user ADD INDEX idx_user_father_id (father_id asc);
ALTER TABLE user ADD INDEX idx_user_root_father_id (root_father_id asc);
ALTER TABLE user ADD INDEX idx_user_root_invite_id (invite_id asc);
ALTER TABLE user ADD INDEX idx_user_root_wallet_id (wallet_id asc);

ALTER TABLE phone ADD INDEX idx_phone_user_id (user_id asc);
ALTER TABLE nickname ADD INDEX idx_nickname_user_id (user_id asc);
ALTER TABLE header ADD INDEX idx_header_user_id (user_id asc);
ALTER TABLE email ADD INDEX idx_email_user_id (user_id asc);
ALTER TABLE wechat ADD INDEX idx_wechat_user_id (user_id asc);
ALTER TABLE wxrobot ADD INDEX idx_wxrobot_user_id (user_id asc);
ALTER TABLE password ADD INDEX idx_password_user_id (user_id asc);
ALTER TABLE username ADD INDEX idx_username_user_id (user_id asc);
ALTER TABLE secondfa ADD INDEX idx_secondfa_user_id (user_id asc);
ALTER TABLE address ADD INDEX idx_address_user_id (user_id asc);
ALTER TABLE idcard ADD INDEX idx_idcard_user_id (user_id asc);
ALTER TABLE company ADD INDEX idx_company_user_id (user_id asc);
ALTER TABLE homepage ADD INDEX idx_homepage_user_id (user_id asc);

ALTER TABLE phone ADD INDEX idx_phone_phone (phone asc);
ALTER TABLE email ADD INDEX idx_email_email (email asc);
ALTER TABLE username ADD INDEX idx_username_username (username asc);
ALTER TABLE wechat ADD INDEX idx_wechat_open_id (open_id asc);
ALTER TABLE wechat ADD INDEX idx_wechat_union_id (union_id asc);
ALTER TABLE wechat ADD INDEX idx_wechat_fuwuhao (fuwuhao asc);

ALTER TABLE pay ADD INDEX idx_pay_user_id (user_id asc);
ALTER TABLE pay ADD INDEX idx_pay_wallet_id (wallet_id asc);
ALTER TABLE pay ADD CONSTRAINT idx_pay_pay_id UNIQUE (pay_id);

ALTER TABLE back ADD INDEX idx_back_user_id (user_id asc);
ALTER TABLE back ADD INDEX idx_back_wallet_id (wallet_id asc);
ALTER TABLE back ADD CONSTRAINT idx_back_back_id UNIQUE (back_id);

ALTER TABLE withdraw ADD INDEX idx_withdraw_user_id (user_id asc);
ALTER TABLE withdraw ADD INDEX idx_withdraw_wallet_id (wallet_id asc);
ALTER TABLE withdraw ADD CONSTRAINT idx_withdraw_withdraw_id UNIQUE (withdraw_id);

ALTER TABLE defray ADD INDEX idx_defray_user_id (user_id asc);
ALTER TABLE defray ADD INDEX idx_defray_wallet_id (wallet_id asc);
ALTER TABLE defray ADD CONSTRAINT idx_defray_defray_id UNIQUE (defray_id);

ALTER TABLE invoice ADD INDEX idx_invoice_user_id (user_id asc);
ALTER TABLE invoice ADD INDEX idx_invoice_wallet_id (wallet_id asc);
ALTER TABLE invoice ADD CONSTRAINT idx_invoice_invoice_id UNIQUE (invoice_id);

ALTER TABLE wallet_record ADD INDEX idx_invoice_user_id (user_id asc);
ALTER TABLE wallet_record ADD INDEX idx_invoice_wallet_id (wallet_id asc);
ALTER TABLE wallet_record ADD INDEX idx_invoice_funding_id (funding_id asc);

ALTER TABLE message ADD INDEX idx_message_user_id (user_id asc);
ALTER TABLE sms_message ADD INDEX idx_sms_message_phone (phone asc);
ALTER TABLE email_message ADD INDEX idx_email_message_email (email asc);
ALTER TABLE fuwuhao_message ADD INDEX idx_fuwuhao_message_open_id (open_id asc);
ALTER TABLE wxrobot_message ADD INDEX idx_wxrobot_message_webhook (webhook asc);
ALTER TABLE audit ADD INDEX idx_audit_user_id (user_id asc);

ALTER TABLE uncle ADD INDEX idx_uncle_user_id (user_id asc);
ALTER TABLE uncle ADD INDEX idx_uncle_uncle_id (uncle_id asc);

ALTER TABLE work_order ADD INDEX idx_work_order_uid (uid asc);
ALTER TABLE work_order ADD INDEX idx_work_order_user_id (user_id asc);

ALTER TABLE work_order_communicate ADD INDEX idx_work_order_communicate_order_id (order_id asc);

ALTER TABLE work_order_communicate_file ADD INDEX idx_work_order_communicate_file_key (`key` asc);

ALTER TABLE discount_buy ADD INDEX idx_discount_buy_user_id (user_id asc, discount_id asc, year asc, month asc, days asc);

ALTER TABLE coupons ADD INDEX idx_coupons_user_id (user_id asc);

ALTER TABLE oauth2_record ADD INDEX idx_oauth2_record_user_id (user_id asc, web_id asc);

ALTER TABLE oauth2_baned ADD INDEX idx_oauth2_baned_user_id (user_id asc, web_id asc);

ALTER TABLE token_record ADD INDEX idx_token_record_create_at (create_at asc);

ALTER TABLE login_controller ADD INDEX idx_login_controller_user_id (user_id asc);

ALTER TABLE face_check ADD INDEX idx_face_check_check_id (check_id asc);

ALTER TABLE access_record ADD INDEX idx_access_record_requests_id_prefix (request_id_prefix asc);
ALTER TABLE access_record ADD INDEX idx_access_record_create_at (create_at asc);

ALTER TABLE oss_file ADD INDEX idx_oss_file_fid (fid asc);

ALTER TABLE website_funding ADD INDEX idx_website_funding_web_id (web_id asc);
ALTER TABLE website_funding ADD INDEX idx_website_funding_funding_id (funding_id asc);
