# nginx_qps
nginx qps watch

## 功能

查看当前nginx总体QPS

编译main.go 生成 stubstatus(.exe)可执行文件,或者直接使用编译好的文件

## nginx 配置
```
nginx -V
查看configure arguments是否有
--with-http_stub_status_module
```
新增server:
```
server{
   listen 9100;# 按需修改
   server_name _;# 按需修改
   location / {
         stub_status on;# 必须
         access_log off;
    }
}
```

重启后运行即可,如
```
linux在可运行文件目录运行:./stubstatus -url http://127.0.0.1:9100
windows在可运行文件目录运行：stubstatus.exe -url http://127.0.0.1:9100
不添加url参数默认请求地址：http://127.0.0.1:9100
如果使用其他地址：执行时加上参数如：stubstatus(.exe) -url http://127.0.0.2/test
```

结果每秒输出一次，如下：
```
时间:2019-07-19 17:09:54
当前QPS: 1
当前连接数: 1
最大QPS:7 发生时间:2019-07-19 17:09:36
最大连接数:1 发生时间:2019-07-19 17:08:24
```