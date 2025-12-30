# 开发文档

本文档介绍如何搭建开发环境、代码规范和贡献指南。

## 📋 目录

- [开发环境搭建](#开发环境搭建)
- [项目结构](#项目结构)
- [代码规范](#代码规范)
- [开发流程](#开发流程)
- [贡献指南](#贡献指南)

## 开发环境搭建

### 前置要求

- **Go**: 1.21+
- **Node.js**: 18+
- **Python**: 3.10+ (用于Python SDK开发)
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **PostgreSQL**: 15+ (或使用Docker)
- **Redis**: 7+ (或使用Docker)

### 后端开发环境

1. 克隆项目：

```bash
git clone <repository-url>
cd blog
```

2. 配置后端：

```bash
cd backend
cp config/config.example.yaml config/config.yaml
# 编辑 config.yaml 配置数据库和Redis连接
```

3. 安装依赖：

```bash
go mod download
```

4. 启动开发环境：

```bash
# 从项目根目录
./scripts/start-dev.sh
```

5. 运行后端（开发模式）：

```bash
cd backend
go run cmd/server/main.go
```

后端服务运行在 `http://localhost:8080`

### 前端开发环境

1. 安装依赖：

```bash
cd frontend
npm install
```

2. 启动开发服务器：

```bash
npm run dev
```

前端服务运行在 `http://localhost:5173`

### Python SDK开发环境

1. 安装uv（推荐）：

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

2. 安装依赖：

```bash
cd python-sdk
uv pip install -e ".[dev]"
```

3. 运行示例：

```bash
python examples/task_reporter_example.py
```

### 数据库迁移

数据库迁移文件位于 `backend/migrations/` 目录。

迁移顺序：
1. `001_init_schema.sql` - 初始化表结构
2. `002_add_indexes.sql` - 添加索引
3. `003_add_fulltext_search.sql` - 全文搜索
4. `004_add_triggers.sql` - 触发器
5. `005_add_init_data.sql` - 初始数据
6. `006_add_performance_indexes.sql` - 性能索引

在Docker环境中，迁移会自动执行。本地开发需要手动执行：

```bash
docker-compose exec postgres psql -U blog_user -d blog_db -f /path/to/migration.sql
```

## 项目结构

### 后端结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # 入口文件
├── internal/
│   ├── config/              # 配置管理
│   ├── handler/             # HTTP处理器
│   ├── middleware/          # 中间件
│   ├── models/              # 数据模型
│   ├── repository/          # 数据访问层
│   ├── service/             # 业务逻辑层
│   ├── router/              # 路由配置
│   ├── scheduler/           # 定时任务
│   ├── websocket/           # WebSocket
│   └── pkg/                 # 工具包
│       ├── crypto/          # 加密工具
│       ├── db/             # 数据库工具
│       ├── jwt/            # JWT工具
│       ├── logger/         # 日志工具
│       ├── redis/          # Redis工具
│       └── response/       # 响应格式化
├── config/                  # 配置文件
├── migrations/              # 数据库迁移
└── docs/                    # Swagger文档
```

### 前端结构

```
frontend/
├── src/
│   ├── api/                 # API接口
│   ├── assets/              # 静态资源
│   ├── components/          # 组件
│   │   ├── article/         # 文章组件
│   │   ├── common/          # 通用组件
│   │   ├── editor/          # 编辑器组件
│   │   ├── fingerprint/     # 指纹组件
│   │   ├── layout/         # 布局组件
│   │   └── upload/          # 上传组件
│   ├── composables/         # 组合式函数
│   ├── locales/             # 国际化
│   ├── router/              # 路由配置
│   ├── store/               # 状态管理
│   ├── utils/               # 工具函数
│   └── views/               # 页面视图
├── public/                  # 公共资源
└── index.html               # HTML模板
```

### Python SDK结构

```
python-sdk/
├── src/
│   └── blog_sdk/
│       ├── crawler/         # 爬虫模块
│       ├── monitor/         # 监控模块
│       └── utils/           # 工具模块
├── examples/                # 示例代码
└── tests/                   # 测试代码
```

## 代码规范

### Go代码规范

1. **命名规范**
   - 包名：小写，简短
   - 函数名：驼峰命名，公开函数首字母大写
   - 变量名：驼峰命名
   - 常量名：全大写，下划线分隔

2. **注释规范**
   - 公开函数必须有注释
   - 注释以函数名开头
   - 使用 `//` 进行单行注释

3. **错误处理**
   - 所有错误必须处理
   - 使用 `errors.New()` 创建错误
   - 使用 `fmt.Errorf()` 包装错误

4. **代码格式**
   - 使用 `gofmt` 格式化代码
   - 使用 `golangci-lint` 检查代码

示例：

```go
// GetUserByID 根据ID获取用户
func GetUserByID(id uint) (*User, error) {
    if id == 0 {
        return nil, errors.New("invalid user id")
    }
    // ...
}
```

### Vue代码规范

1. **组件命名**
   - 组件名使用PascalCase
   - 文件名与组件名一致

2. **代码组织**
   - `<script setup>` 优先
   - 使用组合式API
   - 按功能组织代码

3. **样式规范**
   - 使用Less预处理器
   - 使用CSS变量
   - 遵循BEM命名规范（可选）

4. **代码格式**
   - 使用Prettier格式化
   - 使用ESLint检查

示例：

```vue
<template>
  <div class="article-list">
    <ArticleItem v-for="article in articles" :key="article.id" :article="article" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import ArticleItem from '@/components/article/ArticleItem.vue'

const articles = ref([])

onMounted(() => {
  // 加载文章列表
})
</script>

<style lang="less" scoped>
.article-list {
  padding: 20px;
}
</style>
```

### Python代码规范

1. **命名规范**
   - 函数名：snake_case
   - 类名：PascalCase
   - 常量名：UPPER_SNAKE_CASE

2. **类型提示**
   - 使用类型提示
   - 使用 `typing` 模块

3. **文档字符串**
   - 使用docstring
   - 遵循Google风格

4. **代码格式**
   - 使用Black格式化
   - 使用Ruff检查

示例：

```python
from typing import Optional, Dict, Any

class TaskReporter:
    """Task Reporter for reporting crawler task status."""
    
    def register_task(
        self,
        task_id: str,
        task_name: str,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Register a new crawler task.
        
        Args:
            task_id: Unique task identifier
            task_name: Task name
            metadata: Optional metadata dictionary
            
        Returns:
            Task data from API response
        """
        # ...
```

## 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 开发功能

- 编写代码
- 编写测试（可选）
- 更新文档

### 3. 提交代码

遵循Conventional Commits规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/配置

示例：

```bash
git commit -m "feat(article): 添加文章搜索功能

- 实现全文搜索
- 添加搜索高亮
- 更新API文档"
```

### 4. 推送代码

```bash
git push origin feature/your-feature-name
```

### 5. 创建Pull Request

在GitHub上创建PR，等待代码审查。

## 贡献指南

### 如何贡献

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

### 代码审查标准

- 代码符合规范
- 有适当的注释
- 有测试覆盖（可选）
- 有文档更新
- 无安全隐患
- 无性能问题

### 报告问题

在GitHub Issues中报告问题，包括：
- 问题描述
- 复现步骤
- 预期行为
- 实际行为
- 环境信息

### 功能请求

在GitHub Issues中提出功能请求，包括：
- 功能描述
- 使用场景
- 预期效果

## 开发工具

### 后端开发

- **IDE**: GoLand / VS Code
- **调试**: Delve
- **测试**: Go test
- **文档**: Swagger

### 前端开发

- **IDE**: VS Code
- **调试**: Vue DevTools
- **测试**: Vitest (可选)
- **构建**: Vite

### Python SDK开发

- **IDE**: PyCharm / VS Code
- **调试**: Python debugger
- **测试**: pytest
- **格式化**: Black
- **检查**: Ruff

## 常见问题

### Q: 如何调试后端代码？

A: 使用Delve调试器：

```bash
dlv debug cmd/server/main.go
```

### Q: 如何查看Swagger文档？

A: 启动服务后访问 `http://localhost:8080/swagger/index.html`

### Q: 如何添加新的API接口？

A: 
1. 在 `handler/` 中创建处理器
2. 在 `service/` 中实现业务逻辑
3. 在 `router/` 中注册路由
4. 添加Swagger注释

### Q: 如何添加新的前端页面？

A:
1. 在 `views/` 中创建页面组件
2. 在 `router/index.js` 中添加路由
3. 在 `locales/` 中添加国际化文本

### Q: 如何运行数据库迁移？

A: 在Docker环境中，迁移会自动执行。本地开发需要手动执行SQL文件。

## 获取帮助

- 查看项目文档
- 查看GitHub Issues
- 联系维护者

