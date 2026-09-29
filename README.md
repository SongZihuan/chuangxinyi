# 创信易

企业客户统一身份认证与预付金支付平台。同一企业旗下的其他业务平台接入本站，由本站统一完成身份登记、登录注册、收款与退款。

线上体验：<https://auth.wuntsong.com>

## 能力

**用户认证**　支持手机号、邮箱、身份证、企业营业执照认证。业务平台把用户跳转到创信易完成登录或注册，创信易为其分配 OpenID 等信息并回调给业务平台。

**收款退款**　支持微信支付、支付宝支付，也支持人工收费录入，并可退款。业务平台创建订单后跳转至创信易支付，用户可选择微信支付、支付宝或预付款余额。

**接入站点**　接入方可查看自己的接入信息，以及用户支付所产生的收益。

## 仓库结构

单仓（monorepo）内含两个可独立构建的应用，无子模块：

```
backend/    Go 后端（go.mod 在此目录，勿上提到仓库根）
console/    Vue3 前端
```

## 技术栈

| | 后端 | 前端 |
|---|---|---|
| 语言/框架 | Go 1.21、go-zero（API 定义在 `api/v1`）、GORM | Vite 4、Vue 3、Pinia、Vue Router |
| 存储 | MySQL 8.0+、Redis | — |
| UI | — | Element Plus |

支持启动多个后端实例，实例间通过 `user.group` 相互通信同步数据，可实现冗余容灾。建议容器化运行（Docker / K8s）。

## 快速开始

### 后端

配置目录默认为启动目录下的 `etc/`，可用 `-f` 指定其他路径；该目录已在 `.gitignore` 中，需自行创建 `etc/config.yaml`。全部配置项的逐项说明见 [`backend/README.md`](backend/README.md)。

```bash
cd backend
go run ./src/cmd/user          # 用户中心，默认监听 3350
go run ./src/cmd/sqlclear      # 数据库定期清理（独立进程）
```

容器部署时把配置目录挂到 `/usr/local/share/backend/etc`。

### 前端

前端脚本依赖运行时的工作目录（`vite.config.ts` 以 `process.cwd()` 解析 root 与 env），必须在 `console/` 下执行。

```bash
cd console
pnpm install
pnpm dev                       # http://localhost:8080，/adminApi 代理到 http://localhost:3350/api/v1/
```

需要 Node ≥ 16。私有环境（内网、测试服）的地址请写进 `.env.[mode].local`，该形式已被 `.gitignore` 屏蔽，不会入库。

### 镜像构建

Dockerfile 内的 `COPY` 以**各自子目录**为构建上下文，因此命令必须带 `-f` 并把上下文指向子目录；用仓库根当上下文会因根目录没有 `go.mod` 而失败。

```bash
docker build -f backend/user.dockerfile backend
docker build -f console/nginx.dockerfile console
```

测试环境镜像改用各自的 `*.test.dockerfile`。

## 协议

[MIT](LICENSE)
