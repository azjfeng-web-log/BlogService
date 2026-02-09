# BlogService

基于 Go + Gin + GORM 的博客后端服务。

## 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.21 |
| Web框架 | Gin | 1.9.1 |
| ORM | GORM | 1.25.5 |
| 数据库 | MySQL | - |
| 认证 | JWT | 5.2.0 |
| 加密 | bcrypt | - |

## 项目结构

```
BlogService/
├── main.go                 # 入口 + 路由配置
├── go.mod
└── internal/
    ├── handler/            # HTTP 处理器
    │   ├── article.go      # 文章 API
    │   ├── comment.go      # 评论 API
    │   └── user.go         # 用户 API
    ├── middleware/         # 中间件
    │   └── middleware.go   # CORS + JWT 认证
    ├── model/              # 数据模型
    │   ├── db.go           # 数据库初始化
    │   ├── user.go         # 用户模型
    │   ├── article.go      # 文章/分类/标签/点赞/收藏模型
    │   └── comment.go      # 评论模型
    └── util/               # 工具
        ├── jwt.go          # JWT 生成/解析
        └── response.go     # 统一响应格式
```

## 数据模型

| 模型 | 说明 |
|------|------|
| User | 用户 (用户名/密码/昵称/头像/角色) |
| Article | 文章 (标题/摘要/内容/封面/分类/标签/浏览量/点赞数/评论数) |
| Category | 分类 |
| Tag | 标签 |
| Comment | 评论 (支持嵌套回复) |
| ArticleLike | 点赞记录 |
| ArticleCollect | 收藏记录 |

## API 接口

服务端口: `6000`，前缀: `/api_v2`

### 认证

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | /auth/login | 登录 |
| POST | /auth/register | 注册 |

### 用户 (需认证)

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | /user/info | 获取用户信息 |
| PUT | /user/info | 更新用户信息 |
| PUT | /user/password | 修改密码 |

### 文章

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | /article/list | 文章列表 (分页/筛选) | 否 |
| GET | /article/recommend | 推荐文章 (Top5) | 否 |
| GET | /article/archive | 文章归档 | 否 |
| GET | /article/:id | 文章详情 | 否 |
| GET | /article/:id/interaction | 点赞/收藏状态 | 是 |
| POST | /article/:id/like | 点赞/取消 | 是 |
| POST | /article/:id/collect | 收藏/取消 | 是 |
| POST | /article | 创建文章 | 是 |
| PUT | /article/:id | 更新文章 | 是 |
| DELETE | /article/:id | 删除文章 | 是 |

### 分类/标签

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | /category/list | 分类列表 |
| GET | /tag/list | 标签列表 |

### 评论

| 方法 | 端点 | 说明 | 认证 |
|------|------|------|------|
| GET | /comment/list | 评论列表 | 否 |
| POST | /comment/add | 添加评论 | 是 |
| DELETE | /comment/:id | 删除评论 | 是 |

## 快速开始

1. 创建数据库
```bash
mysql -u root -p -e "CREATE DATABASE blog DEFAULT CHARACTER SET utf8mb4;"
```

2. 修改数据库配置 (`internal/model/db.go`)
```go
dsn := "root:password@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
```

3. 运行
```bash
go mod tidy
go run main.go
```

## 响应格式

```json
// 成功
{"code": 0, "message": "success", "data": ...}

// 失败
{"code": 错误码, "message": "错误信息"}
```
