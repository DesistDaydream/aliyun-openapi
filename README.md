# aliyun-openapi

一个调用阿里云 OpenAPI 的工具仓库

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

全部删除后逐一添加

```go
go run main.go alidns -u ${用户名} -d desistdaydream.ltd -o add -f desistdaydream.ltd.xlsx
```

批量添加

```go
go run main.go alidns -u ${用户名} -o batch -O RR_ADD -d desistdaydream.ltd -f /mnt/d/Downloads/desistdaydream.ltd.xlsx
```

批量删除

```go
go run main.go alidns -u ${用户名} -o batch -O RR_DEL -d desistdaydream.ltd -f /mnt/d/Downloads/desistdaydream.ltd.xlsx
```
