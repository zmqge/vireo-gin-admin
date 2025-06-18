# Vireo Gin Admin 项目介绍

## 项目概述
Vireo Gin Admin 是一个基于Go+Gin+Gorm开发的后台管理系统，提供完整的权限管理、用户管理、通知公告等功能模块。

## 功能特性

- :scroll: 优雅实现 `RESTful API`，采用接口化编程范式，让您的 API 设计更加专业规范
- :house: 采用清晰简洁的模块化架构，让代码结构一目了然，维护升级更轻松自如
- :rocket: 基于高性能 `GIN` 框架，集成丰富实用的中间件（身份认证、跨域、日志、权限控制、容错等），助您快速构建企业级应用
- :closed_lock_with_key: 基于RBAC的用户权限控制和基于部门的数据权限控制，让安全防护固若金汤
- :page_facing_up: 基于功能强大的 `GORM 2.0` ORM 框架，优雅处理数据库操作，大幅提升开发效率
- :memo: 基于高性能 `Zap` 日志框架，配合 Context 链路追踪，让系统运行状态清晰透明，问题排查无所遁形
- :key: 整合久经考验的 `JWT` 认证机制，让用户身份验证更加安全可靠
- :100: 采用无状态设计，支持水平扩展，搭配 Redis 实现动态权限管理，让您的系统轻松应对高并发
- :rocket: 注解路由和注解权限，采用编译时注解处理，零运行时开销设计，让代码运行时开销最小，提升性能

## 在线演示

[https://admin.zmqiang.com](https://admin.zmqiang.com)

## 项目仓库
- 后端:
  - [Gitee](https://gitee.com/zmqge/vireo-gin-admin)
  - [GitHub](https://github.com/zmqge/vireo-gin-admin)
- 前端:
  - [Gitee](https://gitee.com/zmqge/vue3-gin-gorm)
  - [GitHub](https://github.com/zmqge/vue3-gin-gorm)

## 技术栈
- **后端框架**: Gin
- **数据库**: MySQL
- **ORM**: GORM
- **缓存**: Redis
- **认证**: JWT

## 项目结构
```
app/
|--  admin/
  |--  controllers/   # 控制器层
    models/        # 数据模型层
    repositories/  # 数据访问层
    services/      # 业务逻辑层

cmd/
  generator/      # 代码生成器
  permgen/        # 权限生成器
  routegen/       # 路由生成器

config/          # 配置文件

pkg/
  auth/           # 认证相关
  cache/          # 缓存处理
  database/       # 数据库配置
  logger/         # 日志处理
  middleware/     # 中间件
  redis/          # Redis客户端
  response/       # 统一响应处理
  scopes/         # 数据范围处理

routes/          # 路由定义
utils/           # 工具函数

vendor/          # 第三方依赖
```
## 快速开始

1. 克隆项目:
```bash
git clone https://gitee.com/zmqge/vireo-gin-admin.git
```

2. 安装依赖:
```bash
go mod tidy
```

3. 配置数据库:
- 修改config/config.yaml中的数据库配置
- 导入db/dump-vireo_gin_admin-202506161508.sql

4. 运行项目:
```bash
go run main.go
```

## 项目仓库
- 后端:
  - [Gitee](https://gitee.com/zmqge/vireo-gin-admin)
  - [GitHub](https://github.com/zmqge/vireo-gin-admin)
- 前端:
  - [Gitee](https://gitee.com/zmqge/vue3-gin-gorm)
  - [GitHub](https://github.com/zmqge/vue3-gin-gorm)

## 功能详细说明

### 登录设备记录
- 记录用户登录IP地址
- 记录登录时间
- 解析User-Agent获取设备、操作系统和浏览器信息

### 公告功能
- 支持多接收者
- 公告状态管理（已读/未读）
- 公告分类

## 文档
- [项目文档](https://zmqiang.com)
## 贡献指南
欢迎提交Pull Request，请确保代码风格一致并通过测试。