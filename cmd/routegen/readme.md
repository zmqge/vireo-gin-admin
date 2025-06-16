# 路由生成器使用说明

## 功能概述

路由生成器是一个基于代码注解的自动化工具，用于扫描控制器中的路由注解并自动生成Gin框架的路由注册代码。主要功能包括：

1. 自动扫描控制器目录中的路由注解
2. 生成标准化的路由注册代码
3. 自动处理控制器依赖注入
4. 支持中间件自动装配（JWT、RBAC等）
5. 支持路由分组

## 使用方法

### 1. 添加路由注解

在控制器方法上添加 `@Route` 注解，示例：

```go
// @Route(method="GET", path="/api/users", middlewares=["jwt", "rbac"])
func (c *UserController) List(ctx *gin.Context) {
    // 控制器逻辑
}
```

注解参数说明：
- `method`: HTTP方法（GET/POST/PUT/DELETE等）
- `path`: 路由路径
- `middlewares`: 中间件列表（可选）

### 2. 配置扫描目录

在 `config.yaml` 中配置控制器扫描目录：

```yaml
controller_dirs:
  - "app/admin"
  - "app/other"
```

### 3. 运行生成器

执行以下命令生成路由代码：

```bash
go run cmd/routegen/main.go
```

## 生成文件说明

生成器会在 `routes/` 目录下创建以下文件：

1. `{模块名}-api.go`: 具体路由注册代码
2. `route.go`: 统一路由入口文件

## 高级功能

### 路由分组

在控制器结构体上添加 `@Group` 注解可实现路由分组：

```go
// @Group(path="/api/users", middlewares=["jwt"])
type UserController struct {
    // 控制器字段
}
```

### 自动依赖注入

生成器会自动分析控制器的构造函数参数，并按需实例化服务和仓储：

```go
// 自动生成的依赖注入代码示例
userRepo := repositories.NewUserRepository(db)
userService := services.NewUserService(userRepo)
userCtrl := controllers.NewUserController(userService)
```

### 中间件自动装配

- 当路由设置了 `permission` 但未声明 `jwt` 中间件时，生成器会自动添加
- RBAC 中间件会根据 `permission` 自动生成
- 其他中间件需显式声明

## 注意事项

1. 控制器文件必须包含有效的 `New{ControllerName}` 构造函数
2. 服务和仓储也需要提供标准构造函数
3. 生成器会先删除所有现有的 `-api.go` 文件再重新生成
4. 路由变更后需要重新运行生成器