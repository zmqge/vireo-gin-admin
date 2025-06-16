# Go MVC 代码生成器使用说明

## 功能概述
本工具用于快速生成基于 Gin 和 GORM 的 MVC 结构代码，包含以下组件：
- Model（模型）
- Controller（控制器）
- Service（服务层）
- Repository（数据访问层）

## 安装与使用

### 1. 安装
```bash
go install github.com/zmqge/vireo-gin-admin/cmd/generator@latest  
```

### 2. 基本使用
```bash
# 基本命令
go run cmd/generator/main.go -entity=Dict -path=app/admin

# 完整参数
go run cmd/generator/main.go -entity=Dict -module=github.com/zmqge/vireo-gin-admin -path=app/admin
```

### 3. 参数说明
| 参数       | 必填 | 说明                          | 示例值                   |
|------------|------|-----------------------------|-------------------------|
| -entity    | 是   | 实体名称（首字母大写）           | Dict                    |
| -module    | 否   | Go模块路径（自动检测当前项目）    | github.com/yourproject  |
| -path      | 否   | 输出目录（默认为当前目录）        | app/admin               |

## 生成的文件结构
```
app/admin/
├── controllers/
│   └── dictController.go
├── models/
│   └── dictModel.go
├── repositories/
│   └── dictRepository.go
└── services/
    └── dictService.go
```

## 功能特性

### 1. 自动生成的API端点
| 方法   | 路径               | 描述           |
|--------|--------------------|----------------|
| GET    | /dicts/:id        | 获取单个字典    |
| GET    | /dicts            | 获取字典列表    |
| POST   | /dicts            | 创建字典        |
| PUT    | /dicts/:id        | 更新字典        |
| DELETE | /dicts/:id        | 删除字典        |

### 2. 自动权限控制
生成的代码包含权限注解，可与权限系统集成：
```go
// @Permission(code="dict:view", desc="查看字典详情")
```

### 3. 统一响应格式
使用 `github.com/zmqge/vireo-gin-admin/pkg/response` 包返回统一格式的响应：

成功响应：
```json
{
  "code": 0,
  "data": {},
  "msg": "success"
}
```

错误响应：
```json
{
  "code": 400,
  "error": "Invalid ID"
}
```

## 自定义模板

如需修改生成模板，可以：

1. 克隆项目代码
2. 修改 `main.go` 中的模板常量
3. 重新编译安装

## 注意事项

1. 确保项目已初始化 Go module (`go mod init`)
2. 生成的代码需要配合以下依赖：
   - `github.com/gin-gonic/gin`
   - `gorm.io/gorm`
   - 项目自有的 `response` 包

## 示例输出

生成 `Dict` 实体后的控制器示例：
```go
func (c *DictController) getDict(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        response.BadRequest(ctx, "Invalid ID")
        return
    }
    
    entity, err := c.service.GetDictByID(uint(id))
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx, entity)
}
```