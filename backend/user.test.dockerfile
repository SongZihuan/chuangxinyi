FROM registry.cn-guangzhou.aliyuncs.com/wuntsong_pub/ngolang:2023-11-18-16-24-07 as builder

WORKDIR /tmp/backend
RUN mkdir -p /tmp/backend

COPY . .
RUN go mod tidy
RUN go build -o backend -ldflags="all=-w -s" gitee.com/wuntsong-auth/backend/src/cmd/user

FROM registry.cn-guangzhou.aliyuncs.com/wuntsong_pub/nubuntu:2023-11-18-16-20-22

WORKDIR /usr/local/share/backend
RUN mkdir -p /usr/local/share/backend
RUN mkdir -p /usr/local/share/backend/etc

COPY --from=builder /tmp/backend/backend ./
COPY --from=builder /tmp/backend/sql ./sql
COPY --from=builder /tmp/backend/template ./template

RUN groupadd -g 1131 app && useradd -u 1130 -g app app
USER 1130

ENTRYPOINT ["/usr/local/share/backend/backend", "-f", "/usr/local/share/backend/etc"]
