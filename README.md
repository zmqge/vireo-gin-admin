
# vireo-gin-admin

基于Gin框架开发的Go语言后台管理系统

## 功能特性

- 用户认证与授权管理
- 基于RBAC的权限控制
- 部门数据权限控制
- 前后端分离架构
- JWT双Token鉴权机制
- 自动路由生成器
- 权限生成器
- 代码生成器
- 数据库操作封装
- 日志记录与监控

## 技术栈

- 后端: Go (Gin框架)
- 数据库: MySQL/GORM
- 缓存: Redis
- 前端: Vue.js (可选)

## 项目结构

```
├── app/            # 应用核心代码
│   ├── admin/      # 后台管理模块
│   │   ├── controllers/  # 控制器
│   │   ├── models/       # 数据模型
│   │   ├── repositories/ # 数据访问层
│   │   └── services/     # 业务逻辑层
├── cmd/            # 命令行工具
│   ├── routegen/   # 路由生成器
│   └── permgen/    # 权限生成器
├── config/         # 配置文件
├── pkg/           # 公共库
│   ├── auth/      # 认证相关
│   ├── database/  # 数据库操作
│   └── middleware/ # 中间件
├── routes/        # 生成的路由文件
└── main.go        # 应用入口
```

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/zmqge/vireo-gin-admin.git
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
修改`config/config.yaml`中的数据库连接信息

4. 运行应用
```bash
go run main.go
```

## 开发工具

- `cmd/routegen`: 自动生成路由
- `cmd/permgen`: 自动生成权限

## 许可证

MIT License