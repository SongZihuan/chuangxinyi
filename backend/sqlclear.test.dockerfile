FROM registry.cn-guangzhou.aliyuncs.com/wuntsong_pub/ngolang:2023-11-18-16-24-07 as builder

WORKDIR /tmp/sqlclear
RUN mkdir -p /tmp/sqlclear

COPY . .
RUN go mod tidy
RUN go build -o sqlclear -ldflags="all=-w -s" gitee.com/wuntsong-auth/backend/src/cmd/sqlclear

FROM registry.cn-guangzhou.aliyuncs.com/wuntsong_pub/nubuntu:2023-11-18-16-20-22

WORKDIR /usr/local/share/sqlclear
RUN mkdir -p /usr/local/share/sqlclear
RUN mkdir -p /usr/local/share/sqlclear/etc

COPY --from=builder /tmp/sqlclear/sqlclear ./
COPY --from=builder /tmp/sqlclear/sql ./sql
COPY --from=builder /tmp/sqlclear/template ./template

RUN groupadd -g 1131 app && useradd -u 1130 -g app app
USER 1130

ENTRYPOINT ["/usr/local/share/sqlclear/sqlclear", "-f", "/usr/local/share/sqlclear/etc"]
