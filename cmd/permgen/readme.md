# 权限生成器使用说明

## 功能概述
本工具用于从代码注解自动生成权限配置，支持RBAC权限模型。

## 使用方法

### 1. 在代码中添加注解
```go
// @Permission(code="sys:notice:add", name="新建", modules="通知管理", desc="新建通知")
func CreateUser(c *gin.Context) {
    // 业务逻辑
}
```

### 2. 运行生成器
```bash
go run cmd/permgen/main.go
```
提取所有注解权限存入数据库