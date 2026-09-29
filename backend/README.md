# 创信易身份验证系统
## 介绍
创信易是企业客户统一身份认证平台。企业的其他业务平台可接入此身份验证平台，对客户身份进行登记验证。
客户可以使用此平台预充值，在企业的其他业务平台进行消费。

体验：[创新易](https://auth.wuntsong.com)

## 如何启动
推荐docker启动。
将配置文件夹（包含config.yaml）的文件夹挂载到容器的"/usr/local/share/backend/etc"，然后启动即可。
端口默认为3350.

## 推荐配置
关于结构体，参考`src/config`中，以`_model`结尾的文件。
最顶层的结构体位于`config_model.go`中。

```yaml
# config.yaml
mode: develop # 模式 develop或release

user:
  # http的相关设定（其他设定如端口等参考代码中的结构体）
  maxConns: 10000
  maxBytes: 20971520  # 20mb
  
  logServiceName: 创信易用户中心  # 日志名称
  logLevel: info  # 日志等级

  dev:
    metricsUrl: ""

  middlewares:  # 中间件，参考go-zero
    prometheus: true
    maxConns: false
    breaker: false
    shedding: false
    timeout: true
    metrics: true
    maxBytes: true
    gunzip: true

  origin:  # 跨域的origin
    - https://auth.wuntsong.com

  readableName: 创信易用户中心  # 站点名称
  websiteUID: 00000000-0000-0000-0000-000000000000  # 站点UID

  allowMethod:  # 跨域允许的方法
    - GET
    - POST

  allowHeader:  # 跨域允许的请求头
    - Content-Type
    - X-Token
    - X-Domain
    - X-Timestamp
    - X-N
    - X-Sign
    - X-RunMode
    - Origin
    - X-CAPTCHA-Token
    - X-CAPTCHA-Sig
    - X-CAPTCHA-SessionId
    - X-CAPTCHA-Scene
    - X-CAPTCHA-Nvc
    - X-Phone-Token
    - X-Email-Token
    - X-2FA-Token

  url: https://auth.wuntsong.com  # 线上可访问地址，用于生成访问的url

  group: online  # 组名（分布式）

mysql:
  dsn: # {MySQL名称}:{MySQL密码}}@tcp({MySQL地址}:{MySQL端口})/{MySQL数据库}?charset=utf8mb4&parseTime=True&loc=Local
  sqlFilePath: sql  # sql文件的位置

  # 资源配额
  systemResourceLimit: 1000
  workOrderFileLimit: 10
  sonUserLimit: 1000
  uncleUserLimit: 1000
  nephewLimit: 1000
  sameWalletUserLimit: 1000
  sonLevelLimit: 100

  # 数据库清理
  clearCreateAt:
    - tableName: access_record
    - tableName: token_record

  clearDeleteAt:
    - tableName: address
    - tableName: agreement
    - tableName: announcement
    - tableName: application
    - tableName: audit
    - tableName: back
    - tableName: company
    - tableName: coupons
    - tableName: defray
    - tableName: discount
    - tableName: discount_buy
    - tableName: email
    - tableName: email_message
    - tableName: face_check
    - tableName: fuwuhao_message
    - tableName: header
    - tableName: homepage
    - tableName: idcard
    - tableName: invoice
    - tableName: login_controller
    - tableName: menu
    - tableName: message
    - tableName: nickname
    - tableName: oauth2_baned
    - tableName: oauth2_record
    - tableName: oss_file
    - tableName: password
    - tableName: pay
    - tableName: phone
    - tableName: policy
    - tableName: role
    - tableName: secondfa
    - tableName: sms_message
    - tableName: title
    - tableName: uncle
    - tableName: url_path
    - tableName: user
    - tableName: username
    - tableName: wallet
    - tableName: wallet_record
    - tableName: website
    - tableName: website_domain
    - tableName: website_funding
    - tableName: website_ip
    - tableName: website_policy
    - tableName: website_url_path
    - tableName: wechat
    - tableName: withdraw
    - tableName: work_order
    - tableName: work_order_communicate
    - tableName: work_order_communicate_file
    - tableName: wxrobot
    - tableName: wxrobot_message

redis:
  addr: {reids地址}:{redis端口号}
  userName: # reids用户名
  password: # redis密码
  db: 0 # redis数据库

cache:  # 缓存redis
  addr: {reids地址}:{redis端口号}  # 为空标识不启用
  userName: # reids用户名
  password: # redis密码
  db: 0 # redis数据库

aliyun:
  # 阿里云配置
  accessKeyId: 
  accessKeySecret: 
  ocrEndpoint: ocr-api.cn-hangzhou.aliyuncs.com  # 阿里云OCR配置
  # 阿里云短信配置
  smsImportCode:
    sig: 创信易
    template: 
  smsCode:
    sig: 创信易
    template: 
  smsChange:
    sig: 创信易
    template: 
  smsChangePhone:
    sig: 创信易
    template: 
  smsDelete:
    sig: 创信易
    template: 
  smsRegister:
    sig: 创信易
    template: 
  # 身份认证
  identity:
    appKey:
    appSecret: 
    appCode: # 云市场AppCode
    endpoint: https://oss-cn-guangzhou-internal.aliyuncs.com  # 内网可设置内网域名
    bucketName: 
    sign:
      endpoint: https://oss-accelerate.aliyuncs.com  # 内网可设置内网域名
      bucketName: 
  # 验证码
  afs:
    captchaStatus: true
    captchaAppKey: 
    captchaScene: 
    silenceCAPTCHAStatus: true
    silenceCAPTCHAAppKey: 
    silenceCAPTCHAScene: 
  # IP定位
  ip:
    appKey: # 云市场 AppKey
    appSecret: 
    appCode: 
    expiresSecond: 86400 # 24h
  # 头像
  header:
    appKey: # 云市场 AppKey
    appSecret: 
    appCode: 
    endpoint: https://oss-cn-guangzhou-internal.aliyuncs.com  # 内网可设置内网域名
    bucketName: 
    imageStyle: header500x500png
    sign:
      endpoint: https://oss-accelerate.aliyuncs.com  # 必须为公网
      bucketName:
  # 媒体文件
  file:
    endpoint: https://oss-cn-guangzhou-internal.aliyuncs.com  # 内网可设置内网域名
    bucketName: 
    sign:
      endpoint: https://oss-accelerate.aliyuncs.com  # 必须为公网
      bucketName: 
  # OSS工单
  workOrder:
    endpoint: https://oss-cn-guangzhou-internal.aliyuncs.com  # 内网可设置内网域名
    bucketName: 
    sign:
      endpoint: https://oss-accelerate.aliyuncs.com  # 必须为公网
      bucketName: 
  # 发票
  invoice:
    endpoint: https://oss-cn-guangzhou-internal.aliyuncs.com # 内网可设置内网域名
    bucketName: 
    sign:
      endpoint: https://oss-accelerate.aliyuncs.com  # 必须为公网
      bucketName: 

smtp:
  # 邮件服务
  addr: # {smtp地址}:{smtp端口}
  userName: # 账号
  password: # 密码
  fromEmail: # 发件人地址
  sender: # 发件人名字
  templateFilePath: template # 模板地址（不要改）
  sig: # 发件签名

jwt:
  # 各种jwt的配置
  phone:
    subject: PhoneCheck
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  email:
    subject: EmailCheck
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  login:
    subject: Login
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  2fa:
    subject: Login2FA
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  pass2fa:
    subject: Login2FA
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    # 不需要设置过期时间
  user:
    subject: WunTsongUser
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  wechat:
    subject: WeChat
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  idcard:
    subject: IDCard
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  company:
    subject: Company
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  face:
    subject: Face
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  delete:
    subject: Delete
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  check2fa:
    subject: Check2FA
    issuer: GuangZhou WunTsong Information Technology Co,. Ltd
    expiresSecond: 172800
  expireSecond: 172800

cron:
  # cron的执行时间
  permissionUpdate: "@every 30m"
  urlPathUpdate: "@every 30m"
  roleUpdate: "@every 30m"
  menuUpdate: "@every 30m"
  websiteUpdate: "@every 30m"
  websitePermissionUpdate: "@every 30m"
  websiteUrlPathUpdate: "@every 30m"
  applicationUpdate: "@every 30m"

admin:
  adminPhone: # 根管理员手机号

  rootRole:
    roleName: 根管理员
    roleSign: root
    roleDescribe: create by system
    notDelete: true
    notChangeSign: true
    notChangePermissions: true
    resetPermission: true

  userRole:
    roleName: 优质用户
    roleSign: normaluser
    roleDescribe: create by system
    notDelete: true
    notChangeSign: true
    notChangePermissions: false
    resetPermission: true

  anonymousRole:
    roleName: 陌生用户
    roleSign: anonymoususer
    roleDescribe: create by system
    notDelete: true
    notChangeSign: true
    notChangePermissions: false
    resetPermission: true

  menuDepth: 2

  icp1: # ICP备案号
  icp2: # ICP许可证号
  gongan:  # 公安备案号
  copyright: # 版权

wechat:
  # 微信开放平台网页登录的App配置
  appID: # appID
  appSecret: # 私钥直接粘贴上来

totp:
  # totp的配置
  issuerName: auth.wuntsong.com

alipay:
  # 支付宝配置
  sandbox: false
  appID: # appID
  privateKey: # 私钥直接粘贴上来
  publicCert: etc/appPublicCert.crt  # App证书地址
  alipayPublicCert: etc/alipayPublicCert.crt  # 支付宝证书地址
  alipayRootCert: etc/alipayRootCert.crt  # 支付吧根证书地址
  encryptKey: 
  returnUrl: https://auth.wuntsong.com  # 支付成功回调
  notifyUrl: https://auth.wuntsong.com/api/v1/pay/alipay/noti  # 支付成功通知
  faceReturnUrl: https://auth.wuntsong.com  # 刷脸支付成功通知
  wapQuitUrl: https://auth.wuntsong.com  # wap支付取消通知
  usePCPay: true
  useWapPay: true
  useFaceCheck: true
  useReturnPay: true
  useWithdraw: true

wechatpay:
  appID: # 绑定的AppID，可以用公众号的
  mchID: # 商户号ID
  mchAPIv3Key: 
  publicCert: etc/wechatpayCert.pem  # App证书
  privateKey: etc/wechatpayKey.pem  # App密钥
  returnURL: https://auth.wuntsong.com/api/v1/pay/wechatpay/noti  # 支付成功回调地址
  useNativePay: true
  useH5Pay: true
  useReturnPay: true
  useWithdraw: true

coin:
  id: # 产品ID，不要太长
  name: 桓创信息技术服务预付金  # 产品名字
  price: 1  # 单价
  showUrl: "https://www.wuntsong.com"  # 展示的url
  timeExpireSec: 900  # 15min  支付过期时间
  refundDayLimit: 15  # 可退款时间（消费退款）
  refundReason: 用户主动退款  # 用户退款的原因（消费退款）
  returnDayLimit: 15  # 可退款时间（充值退款）
  withdrawMin: 2000  # 20元 # 提现限额

defray:
  secret: # defray的支付secret aes密钥，32位

sign:
  defrayPriKey: # 发送通知的签名密钥，base64编码

message:
  sender: 创信易小助手  # 站内信名字
  senderLink: ""  # 站内信链接

wxrobot:
  log: # 企业微信推送错误日志webhook
  sender: 创信易小助手  # 企业微信发送名字

fuWuHao:
  appID: # 服务号AppID 
  secret: # 服务号secret
  token: 
  encodingAESKey: 

  # 模板消息
  register:
    templateID: 
    url: https://auth.wuntsong.com

  project:
    templateID: 
    url: https://auth.wuntsong.com

  pay:
    templateID: 
    url: https://auth.wuntsong.com

  oauth2:
    templateID: 
    url: ''  # 不需要设置

  loginSuccess:
    templateID: 
    url: ''  # 不需要设置

  loginFail:
    templateID: 
    url: https://auth.wuntsong.com

  userDelete:
    templateID: 
    url: https://auth.wuntsong.com

  # 服务号菜单设置
  kefu:
    huanChuang: # 企业微信客服 kf 链接
    vxwk: # 企业微信客服 kf 链接

  menu:
    updateMenu: true  # 是否更新菜单，若不需要则关闭
    aboutUsWebsite: https://www.wuntsong.com
    aboutUsContact: |-
      欢迎联系我们
      用户反馈：{反馈邮箱}
      用户举报：{举报邮箱}
      商务联系：{商务邮箱}
      亦可通过客服联系我们，输入“桓创客服”可找到客服入口。
    aboutUsKefu: |-
      若想联系创思域变专属客服，输入：“创思域变客服”
      若想联系其他客服，输入：“桓创客服”
      我们的客户经理将尽快联系您。
    productVXWK: https://www.vxwk.com
    productWallpaper: https://www.wuntsong.com/productinfo/<商品ID>.html
    auth: https://auth.wuntsong.com/fuwuhao

password:  # 密码加盐
  salt: 
  frontSalt: 

sqlClear:
  logServiceName: # sql定时清理日志名称


```

## 代码阅读指南
### 入口函数
位于`src/cmd`中，包括：

1. sqlclear: 数据库定期自动清理
2. user: 用户中心后端

### 后端Api
位于`api/v1`文件夹中，使用`go-zero-plus`框架。

### 后端逻辑
位于`src/service/v1/internal/logic`中。

## 关于疑问
关于任何代码的疑问都可以向我们发送issue，我们会尽量回复。
