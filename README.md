# skyline
实时监控日志工具

# 安装
```bash
go get github.com/gfreezy/skyline
```

# 使用方法
```bash
skyline -h

Usage: skyline [-d] -c <config> [-r <regexp>]

Monitor logs

Test regexp: skyline -r "regexp to test" "text to be tested"

Options:
    -d, --debug             debug
    -c, --config=<config>   config file
    -r, --regexp=<regexp>   test regexp
    -h, --help              usage (-h) / detailed help text (--help)

```

# 测试正则表达式
```bash
skyline -r "(\\d+\\.?\\d*)" '223.104.4.148 09/May/2016:13:32:36 +0800 "GET /v2/recipes/lookup.json?version=5.2.3&id=100104998%2C1001022%2C245919%2C36274%2C100190201%2C245754%2C100047523%2C100483677%2C12416%2C100368789%2C100413705%2C9700%2C100562736%2C100372215%2C253410%2C100021480%2C100407002%2C94322%2C100390847%2C100450619%2C1063469%2C100508029%2C100437043%2C21490%2C1040559%2C264993%2C100369718%2C100051667%2C100002522%2C100388240&api_sign=c5c5f6d131bddc40dec10615e66c377f&api_key=0f9f79be1dac5f003e7de6f876b17c00&origin=iphone&sk=BjuxW487QouV_oL1SNpQAA HTTP/1.1" 200 51868 "-" "Mozilla/5.0 (iPhone; CPU iPhone OS 9.2.1 like Mac OS X) xiachufang-iphone/5.2.3 Build/69" 0.167 0.167 192.168.2.147:8902'
```

# 配置文件
```json
{
  "monitors": [
    {
      "log_name_prefix": "nginx",
      "log_file_path": "test.log",
      "filter_items": [
        {
          "item_name_prefix": "time1",
          "cycle": 2,
          "match_str": "(\\d+\\.?\\d*)",
          "threshold": 1
        },
        {
          "item_name_prefix": "time10",
          "cycle": 5,
          "match_str": "(\\d+\\.?\\d*)",
          "threshold": 10
        }
      ]
    },
    {
      "log_name_prefix": "nginx2",
      "log_file_path": "test2.log",
      "filter_items": [
        {
          "item_name_prefix": "time1",
          "cycle": 1,
          "match_str": "(\\d+\\.?\\d*)",
          "threshold": 1
        },
        {
          "item_name_prefix": "time10",
          "cycle": 3,
          "match_str": "(\\d+\\.?\\d*)",
          "threshold": 10
        }
      ]
    }
  ],
  "warnings": [
    {
      "formula": "nginx_time1_cnt > 20",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过1秒的请求数超过20个",
      "alert_command": "nc 127.0.0.1 8888"
    },
    {
      "formula": "nginx_time1_avg > 1.1",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过1秒的请求数的平均响应时间超过1.1秒",
      "alert_command": ""
    },
    {
      "formula": "nginx_time10_cnt > 20",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过10秒的请求数超过20个",
      "alert_command": ""
    },
    {
      "formula": "nginx_time10_avg > 1.1",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过10秒的请求数的平均响应时间超过1.1秒"
    },
    {
      "formula": "nginx2_time1_cnt > 20",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过1秒的请求数超过20个"
    },
    {
      "formula": "nginx2_time1_avg > 1.1",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过1秒的请求数的平均响应时间超过1.1秒"
    },
    {
      "formula": "nginx2_time10_cnt > 20",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过10秒的请求数超过20个"
    },
    {
      "formula": "nginx2_time10_avg > 1.1",
      "warning_filter": "1/1",
      "alert_name": "5秒内请求时间超过10秒的请求数的平均响应时间超过1.1秒"
    }
  ]
}
```