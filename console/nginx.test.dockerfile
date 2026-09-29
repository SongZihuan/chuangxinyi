FROM registry.cn-guangzhou.aliyuncs.com/wuntsong_pub/nnode:2023-12-10-10-11-38 as builder

WORKDIR /tmp/build
RUN mkdir -p /tmp/build

COPY . .

# 编译项目
RUN pnpm install
RUN pnpm run build:test

FROM registry.cn-guangzhou.aliyuncs.com/wuntsong_pub/nnginx:2023-11-29-22-58-29

USER 0

WORKDIR /var/www
RUN mkdir -p /var/www/project

COPY --from=builder /tmp/build/dist /var/www/project
COPY ./nginx/* /etc/nginx/conf.d

USER 101
