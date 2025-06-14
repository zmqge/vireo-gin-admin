以下是整合了 `authController.go` 和 `userController.go` 的完整 API 文档的 Markdown 格式：

---

# API 文档

## 认证接口

### POST /api/v1/auth/login

用户登录

**Parameters**

| Name        | In   | Type   | Required | Description      |
|-------------|------|--------|----------|------------------|
| username    | body | string | true     | 用户名           |
| password    | body | string | true     | 密码             |
| captchaKey  | body | string | true     | 验证码 ID        |
| captchaCode | body | string | true     | 验证码答案       |

**Responses**

| Code | Description          | Schema                          |
|------|----------------------|---------------------------------|
| 200  | 登录成功             | `{ "tokenType": "Bearer", "accessToken": "string", "refreshToken": "string", "expiresIn": 3600 }` |
| 400  | 无效的请求参数       | `{ "code": "400", "msg": "无效的请求参数" }` |
| 401  | 用户名或密码错误     | `{ "code": "401", "msg": "用户名或密码错误" }` |
| 500  | 服务器内部错误       | `{ "code": "500", "msg": "服务器内部错误" }` |

---

### GET /api/v1/auth/captcha

获取验证码

**Responses**

| Code | Description          | Schema                          |
|------|----------------------|---------------------------------|
| 200  | 获取验证码成功       | `{ "captchaKey": "string", "captchaBase64": "string" }` |
| 500  | 生成验证码失败       | `{ "code": "99999", "msg": "生成验证码失败" }` |

---

### POST /api/v1/auth/refresh

刷新 Access Token

**Parameters**

| Name          | In     | Type   | Required | Description      |
|---------------|--------|--------|----------|------------------|
| Refresh-Token | header | string | true     | Refresh Token    |

**Responses**

| Code | Description          | Schema                          |
|------|----------------------|---------------------------------|
| 200  | 刷新成功             | `{ "access_token": "string" }`  |
| 400  | Refresh Token 缺失   | `{ "code": "400", "msg": "Refresh Token 缺失" }` |
| 401  | Refresh Token 无效   | `{ "code": "401", "msg": "Refresh Token 无效" }` |
| 500  | 生成 Access Token 失败 | `{ "code": "500", "msg": "生成 Access Token 失败" }` |

---

### GET /api/v1/auth/profile

获取用户信息

**Responses**

| Code | Description          | Schema                          |
|------|----------------------|---------------------------------|
| 200  | 获取用户信息成功     | `{ "userID": "string", "username": "string" }` |

---

## 用户管理

### GET /api/v1/users/page

获取用户分页列表

**Parameters**

| Name        | In    | Type   | Required | Description                     |
|-------------|-------|--------|----------|---------------------------------|
| keywords    | query | string | false    | 关键字(用户名/昵称/手机号)      |
| status      | query | string | false    | 用户状态                        |
| roleIDs     | query | string | false    | 角色ID                          |
| createTime  | query | string | false    | 创建时间范围                    |
| field       | query | string | false    | 排序字段                        |
| direction   | query | string | false    | 排序方式（正序:ASC；反序:DESC） |
| pageNum     | query | int    | true     | 页码                            |
| pageSize    | query | int    | true     | 每页记录数                      |

**Responses**

| Code | Description | Schema                          |
|------|-------------|---------------------------------|
| 200  | OK          | [PageResultUserPageVO](#pageresultuserpagevo) |
| 400  | Bad Request | [ResultObject](#resultobject)   |

---

### POST /api/v1/users

创建用户

**Parameters**

| Name     | In   | Type   | Required | Description |
|----------|------|--------|----------|-------------|
| username | body | string | true     | 用户名      |
| password | body | string | true     | 密码        |

**Responses**

| Code | Description | Schema                          |
|------|-------------|---------------------------------|
| 200  | OK          | [ResultObject](#resultobject)   |
| 400  | Bad Request | [ResultObject](#resultobject)   |

---

### DELETE /api/v1/users/:id

删除用户

**Parameters**

| Name | In   | Type   | Required | Description |
|------|------|--------|----------|-------------|
| id   | path | string | true     | 用户ID      |

**Responses**

| Code | Description | Schema                          |
|------|-------------|---------------------------------|
| 200  | OK          | [ResultObject](#resultobject)   |
| 400  | Bad Request | [ResultObject](#resultobject)   |

---

### PUT /api/v1/users/:id

更新用户信息

**Parameters**

| Name     | In   | Type   | Required | Description |
|----------|------|--------|----------|-------------|
| id       | path | string | true     | 用户ID      |
| username | body | string | true     | 用户名      |
| status   | body | int    | true     | 用户状态    |

**Responses**

| Code | Description | Schema                          |
|------|-------------|---------------------------------|
| 200  | OK          | [ResultObject](#resultobject)   |
| 400  | Bad Request | [ResultObject](#resultobject)   |

---

### GET /api/v1/users/me

获取当前用户信息

**Responses**

| Code | Description | Schema                          |
|------|-------------|---------------------------------|
| 200  | OK          | [UserInfoVO](#userinfovo)       |
| 401  | Unauthorized | `{ "code": "401", "msg": "用户未认证" }` |

---

## 数据结构

### PageResultUserPageVO

| 字段 | 类型               | 描述         |
|------|--------------------|--------------|
| code | string             | 状态码       |
| data | DataUserPageVO     | 数据         |
| msg  | string             | 消息         |

### DataUserPageVO

| 字段 | 类型               | 描述         |
|------|--------------------|--------------|
| list | []UserPageVO       | 用户列表     |
| total| int64              | 总数         |

### UserPageVO

| 字段       | 类型   | 描述           |
|------------|--------|----------------|
| id         | int64  | 用户ID         |
| username   | string | 用户名         |
| nickname   | string | 昵称           |
| mobile     | string | 手机号         |
| gender     | int    | 性别           |
| avatar     | string | 头像地址       |
| email      | string | 邮箱           |
| status     | int    | 用户状态       |
| createTime | string | 创建时间       |

### ResultObject

| 字段 | 类型   | 描述         |
|------|--------|--------------|
| code | string | 状态码       |
| data | object | 数据         |
| msg  | string | 消息         |

### UserInfoVO

| 字段     | 类型   | 描述         |
|----------|--------|--------------|
| userId   | string | 用户ID       |
| username | string | 用户名       |
| nickname | string | 昵称         |
| avatar   | string | 头像地址     |
| roles    | []string | 角色列表     |
| perms    | []string | 权限列表     |

---

### 说明：
- **认证接口**：包括登录、获取验证码、刷新 Token 和获取用户信息。
- **用户管理接口**：包括获取用户分页列表、创建用户、删除用户、更新用户信息和获取当前用户信息。
- **数据结构**：定义了接口返回的数据结构。

---

### 使用方法：
1. 将上述 Markdown 内容保存为 `api_doc.md` 文件。
2. 在项目文档中引用该文件，或直接发布到文档平台（如 GitBook、Read the Docs 等）。

如果需要进一步扩展或调整，可以根据实际需求修改文档内容。
