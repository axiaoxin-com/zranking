# zranking

![visitor badge](https://visitor-badge.glitch.me/badge?page_id=axiaoxin-com.zranking&left_color=red&right_color=green&left_text=HelloVisitors)


使用 redis zset 实现的排行榜 golang 封装

## 实现原理

使用redis zset，得分相同时，按时间先后进行排序；
将zset score按十进制数拆分，score十进制数字总共固定为16位（超过16位的数会有浮点数精度导致进位的问题），
整数部分用于表示用户排序值val，小数部分表示排行活动结束时间戳（秒）与用户排序值更新时间戳（秒）的差值deltaTs，
小数部分的数字长度由deltaTs的数字长度确定，整数部分最大支持长度则为：16-len(deltaTs)。
比如活动时长为10天，总时间差为864000，长度为6，则deltaTs宽度为6，不够则在前面补0。

## 安装使用

```
go get -u github.com/axiaoxin-com/zranking
```

示例：[_example/main.go](./_example/main.go) 或参考 [zranking_test.go](./zranking_test.go)



## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=axiaoxin-com/zranking&type=Date)](https://star-history.com/#axiaoxin-com/zranking&Date)

