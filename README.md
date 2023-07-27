# aliyun-openapi

一个调用阿里云 OpenAPI 的工具仓库

下面这两个库有一些功能函数非常好用

- util "github.com/alibabacloud-go/tea-utils/v2/service"
- "github.com/alibabacloud-go/tea/tea"

比如：

- util.ToJSONString(err)
- tea.ToMap(err)

这两种都可以让错误信息以标准格式输出，以便提取其中的某些信息。

# 目录结构

cmd

- dns # 处理云解析的工具

pkg

- aliclient # 创建用于访问阿里云 OpenAPI 的实例，调用 API 时都会使用该实例
- alidns # 云解析 API，https://next.api.aliyun.com/document/Alidns
- fileparse # 文件解析。比如解析 alidns 中用到的 Excel
- config # 配置文件解析。比如认证信息等

# 构建项目

```go
go build cmd/dns/dns.go
```

# 云解析常用操作

更新，不存在则创建

```bash
export USER_NAME="用户名"
go run main.go alidns update -u ${USER_NAME} -d superstor.cn -f /mnt/d/tmp/superstor.cn.xlsx
```

全部删除后逐一添加

```bash
go run main.go alidns full-update -u ${USER_NAME} -d superstor.cn -f /mnt/d/tmp/superstor.cn.xlsx
```
