package defaultuid

// DefaultUID 仅限初始化时单协程写入，其余时刻仅限读取
var DefaultUID string
