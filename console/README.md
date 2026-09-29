# 创信易前端（console）

创信易用户中心的 Web 端，包含用户侧登录注册、钱包与充值提现、订单支付，以及接入站点侧的后台管理界面。与后端在同一仓库的 `../backend`，仅通过 HTTP 通信，无目录级依赖。

Vite 4 + Vue 3 + Pinia + Vue Router + Element Plus。

## 开发

```bash
pnpm install
pnpm dev            # 监听 8080
```

`pnpm dev` 会把 `/adminApi` 代理到 `http://localhost:3350/api/v1/`（即本地后端）。要指向别的后端，在 `.env.development.local` 里设 `VITE_PROXY_TARGET`。

**必须在 `console/` 目录下执行脚本**：`vite.config.ts` 用 `process.cwd()` 解析 root 和 env，在仓库根执行会读不到 `.env`，并在 `JSON.parse(env.VITE_OPEN)` 处直接抛错。

要求 Node ≥ 16。

## 构建

| 命令 | 模式 | 读取的环境文件 |
|---|---|---|
| `pnpm build` | 默认（production） | `.env` + `.env.production` |
| `pnpm build:prod` | production | 同上 |
| `pnpm build:test` | test | `.env` + `.env.test` |

产物在 `dist/`。

## 环境变量

`.env` 是所有模式的基线，模式文件（`.env.production`、`.env.test`、`.env.development`）覆盖其中的同名项。仓库里只保留公开的默认地址；内网、测试服等私有地址请写进 `.env.[mode].local` 或 `.env.local`——这类文件已被 `.gitignore` 屏蔽，不会被提交。

## 容器

```bash
docker build -f nginx.dockerfile .
```

构建上下文为本目录。Nginx 配置在 `nginx/project.conf`，监听 8080，走 SPA fallback。
