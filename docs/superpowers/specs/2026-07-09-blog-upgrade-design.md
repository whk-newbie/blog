# 博客系统升级优化设计文档

> 创建日期：2026-07-09
> 状态：设计完成，待进入实施计划

## 1. 概述

本次升级覆盖安全加固、功能增减、UI 改版三个维度，采用**分层推进**策略：基础设施 → 功能 → UI。

## 2. 升级策略：分层推进

```
安全基础设施（加密 + 混淆 + 入口配置）
        ↓
功能层（去爬虫 + AI 翻译 + AI 聊天）
        ↓
UI 层（zyyo.net 风格全站改版）
```

每层独立完成并验证，后续改动天然受前层保护，可独立测试和回滚。

---

## 3. 第一层：安全基础设施

### 3.1 RSA+AES 混合加密通信

**密钥协商流程：**

```
前端                                    后端
  │                                      │
  │──── GET /api/public-key ────────────→│  返回 RSA 公钥 + sessionId
  │←──── { publicKey, sessionId } ──────│
  │                                      │
  │  生成随机 AES-256-GCM 密钥            │
  │  用 RSA 公钥加密 AES 密钥              │
  │                                      │
  │──── POST /api/session/key ──────────→│  { encryptedKey, sessionId }
  │                                      │    RSA 私钥解密 → 得到 AES 密钥
  │                                      │   Redis: sessionId → AES key (TTL 30min)
  │←──── { success } ────────────────────│
  │                                      │
  │════ 后续所有 API 通信 ════════════════│
  │  请求体 AES-GCM 加密                  │
  │  Header: X-Session-Id                │
  │                                      │  从 Redis 取 AES key
  │                                      │  解密请求 → 处理 → 加密响应
```

**后端实现要点：**
- 新增 `GET /api/public-key`：公开接口，返回 RSA 公钥 + sessionId
- 新增 `POST /api/session/key`：协商 AES 密钥，存入 Redis，TTL 30 分钟
- 全局加解密中间件：拦截除白名单外的所有 `/api/*` 请求
- 白名单：`/api/public-key`、`/api/session/key`
- 加解密使用 AES-256-GCM（带认证标签，防篡改）
- RSA 密钥对可配置（配置文件或环境变量），未配置时自动生成

**前端实现要点：**
- Axios 请求拦截器：首次自动协商密钥，之后所有请求自动加解密
- 密钥过期时（收到 401/特定错误码）自动重新协商
- 加密失败时降级提示用户刷新页面

**依赖：**
- Go `crypto/aes`、`crypto/cipher`、`crypto/rsa`（标准库，无需新增依赖）
- 前端 `crypto-js`（已安装）

### 3.2 后台入口可配置

- 利用现有 `system_configs` 表，新增两条配置项：

| 配置键 | 默认值 | 说明 |
|--------|--------|------|
| `admin_path` | `admin` | 后台访问路径 |
| `show_admin_link` | `true` | 是否在前端显示后台入口链接 |

- 前端路由：`/:adminPath/*` 动态匹配后台路由
- 首页/导航栏：根据 `show_admin_link` 决定是否渲染入口链接
- `show_admin_link` 默认 `true`（默认显示），管理员可主动关闭

### 3.3 前端代码高度混淆

- 使用 `rollup-plugin-obfuscator`（Vite 插件）
- 混淆配置：

```javascript
// vite.config.js - build only
obfuscator({
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.75,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.4,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayThreshold: 0.75,
  selfDefending: true,
  debugProtection: true,
  domainLock: ['yourdomain.com'],  // 部署时配置
})
```

- 仅在 `vite build` 时启用（开发模式不混淆）
- 新增 npm 依赖：`rollup-plugin-obfuscator`

---

## 4. 第二层：功能层

### 4.1 删除爬虫监控功能

**删除范围：**
- `python-sdk/` 整个目录
- 后端爬虫相关：handler、service、repository、路由注册
- 前端爬虫监控页面及路由
- 相关数据库迁移文件（如有独立迁移）

**保留范围：**
- 全部指纹收集代码（fingerprint repo/service/handler，前端指纹组件）
- 访问统计代码（visit repo/service/handler）

### 4.2 AI 文章翻译

**数据模型变更（articles 表新增字段）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `title_en` | VARCHAR(500) | 英文标题 |
| `content_en` | TEXT | 英文内容 |
| `summary_en` | VARCHAR(500) | 英文摘要 |
| `is_translated` | BOOLEAN | 是否已翻译，默认 false |
| `translated_at` | TIMESTAMP | 翻译时间 |

**翻译流程：**

```
发布/编辑文章
      │
      ▼
  保存中文内容
      │
      ├── 勾选"翻译为英文"？
      │     │
      │     ├── 是 → 异步调用 AI 翻译 API
      │     │         翻译 title + content + summary
      │     │         保存到 title_en / content_en / summary_en
      │     │         is_translated = true, translated_at = now
      │     │
      │     └── 否 → 跳过
      │
      ▼
  返回结果
```

**关键规则：**
- 翻译是**一次性**的，`is_translated=true` 后不会重复翻译
- 编辑文章时若中文内容变更，`is_translated` 重置为 `false`，提示重新翻译
- 文章列表提供「翻译」按钮，支持对已有文章手动翻译/重新翻译
- 翻译失败不影响文章发布，仅提示翻译失败

**前端语言切换：**
- 导航栏新增中/英切换按钮（🇨🇳 / 🇬🇧）
- 切换语言时，文章列表/详情加载对应语言字段（title vs title_en）
- 语言偏好存入 localStorage
- 若文章未翻译（`is_translated=false`），英文模式下仍显示中文内容并标注"暂未翻译"

### 4.3 AI 提供方管理

**新表 `ai_providers`：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | UUID | 主键 |
| `name` | VARCHAR(100) | 显示名称，如"我的 DeepSeek" |
| `provider_type` | VARCHAR(20) | claude / openai / deepseek / custom |
| `api_key` | TEXT | API Key，AES 加密存储 |
| `base_url` | VARCHAR(500) | API 地址，custom 时必填 |
| `model` | VARCHAR(100) | 模型名称，如 deepseek-chat |
| `is_enabled` | BOOLEAN | 是否启用 |
| `sort_order` | INT | 优先级排序，数字越小越优先 |
| `balance` | DECIMAL | 余额，定期检测更新 |
| `last_check_at` | TIMESTAMP | 最后检测时间 |
| `created_at` | TIMESTAMP | 创建时间 |
| `updated_at` | TIMESTAMP | 更新时间 |

**后台管理页面：**
- AI 提供方列表，支持 CRUD
- 每个提供方：「检测」按钮，调用 API 验证连通性 + 查询余额
- 翻译和聊天时按 `sort_order` 优先选择 `is_enabled=true` 的提供方
- 翻译失败时自动 fallback 到下一个启用的提供方

**支持的内置 Provider：**
- **Claude**：调用 Anthropic Messages API，base_url 默认 `https://api.anthropic.com`
- **OpenAI (ChatGPT)**：调用 OpenAI Chat Completions API，base_url 默认 `https://api.openai.com`
- **DeepSeek**：调用 DeepSeek Chat API，base_url 默认 `https://api.deepseek.com`
- **Custom**：兼容 OpenAI 接口格式的自定义 API，需填写 base_url

### 4.4 AI 聊天框

**位置与交互：**
- 后台右下角悬浮按钮（💬），点击弹出聊天面板
- 类 ChatGPT 对话界面：消息列表 + 输入框 + 发送按钮
- 支持多轮对话（后端维护对话历史）

**写作辅助快捷指令：**
- 输入框上方预设按钮：润色 / 扩写 / 缩写 / 总结 / 翻译
- 当前正在编辑的文章可被带入对话上下文

**技术方案：**
- HTTP POST `/api/ai/chat` + SSE（Server-Sent Events）流式返回
- 比 WebSocket 实现更简单，单向流式足够
- 对话历史由前端暂存，每次请求带上最近 N 轮对话

**新表 `ai_chat_history`（可选，用于持久化对话记录）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | UUID | 主键 |
| `provider_id` | UUID | 使用的 AI 提供方 |
| `role` | VARCHAR(20) | user / assistant |
| `content` | TEXT | 消息内容 |
| `created_at` | TIMESTAMP | 发送时间 |

---

## 5. 第三层：UI 改版

### 5.1 参考风格

参考网站：[zyyo.net](https://zyyo.net/)
参考仓库：[ZYYO666/homepage](https://github.com/ZYYO666/homepage)

**核心风格特征：**
- 极简设计、内容优先
- 侧边栏 + 内容区双栏布局
- 扁平白色卡片 + 极淡阴影 + 小圆角（4px）
- 中文优先排版，系统字体栈
- 标签云彩色点缀
- 统计计数展示
- 克制动效
- 完整暗色模式支持

### 5.2 全局布局

```
┌──────────────┬──────────────────────────────────────┐
│   Sidebar    │          Main Content               │
│   (固定 280px)│          (滚动)                      │
│              │                                     │
│  📷 头像     │  文章卡片 / 文章详情 / 工具 / 时间轴   │
│  名字/简介   │                                     │
│              │                                     │
│  📊 统计     │                                     │
│  文章数      │                                     │
│  运行天数    │                                     │
│              │                                     │
│  🏷️ 标签云  │                                     │
│  #Go #Vue   │                                     │
│  #Python    │                                     │
│              │                                     │
│  🔗 社交链接 │                                     │
│  GitHub     │                                     │
│  Email      │                                     │
└──────────────┴──────────────────────────────────────┘
```

**响应式：**
- ≥768px：双栏布局，侧边栏固定
- <768px：侧边栏收起，顶部汉堡菜单展开

### 5.3 首页

- **侧边栏**：头像、一句话简介、文章数/运行天数统计、标签云、社交图标
- **内容区**：文章卡片流
  - 每张卡片：16:9 封面图 + 标题 + 日期 + 标签彩色徽章
  - 悬停时卡片微升（translateY -2px + 阴影加深）
- **顶部导航栏**（横跨侧边栏上方或独立一行）：
  - 分类链接、搜索图标（点击弹出全屏搜索）、中/英切换、主题切换

### 5.4 时间轴页面（新增）

```
┌──────────────┬──────────────────────────────────────┐
│   Sidebar    │         时间轴                       │
│             │                                     │
│              │  ● 2026年                           │
│              │  │                                  │
│              │  ├── 7月                             │
│              │  │   ├─ 📝 文章标题A  · 07-09       │
│              │  │   └─ 📝 文章标题B  · 07-03       │
│              │  │                                  │
│              │  ├── 6月                             │
│              │  │   └─ 📝 文章标题C  · 06-15       │
│              │  │                                  │
│              │  ● 2025年                           │
│              │  └── 12月                            │
│              │      ├─ 📝 文章标题D  · 12-28       │
│              │      └─ 📝 文章标题E  · 12-10       │
└──────────────┴──────────────────────────────────────┘
```

- 左侧垂直线 + 年份圆点标记
- 按「年 → 月」两级分组
- 每篇文章：标题 + 日期，点击跳转详情
- 按发布时间倒序排列

### 5.5 文章详情页

- 大标题 + 日期/分类/标签元信息
- 正文 Markdown 渲染，自定义样式
- 代码块 highlight.js，匹配整体浅色/暗色配色
- 底部版权声明 + 许可信息
- 侧边栏保留（可提供收起按钮）

### 5.6 工具页

- 保持现有工具功能不变
- 容器改为卡片风格：白色背景 + 极淡阴影 + 4px 圆角
- 表单控件保留 Element Plus（功能需要）

### 5.7 后台管理

- 布局不变（左侧菜单 + 右侧内容区）
- 新增强化：将 Element Plus 组件替换部分为自定义极简风格
- 主色调改为浅灰白背景（`#fafafa`），降低视觉压迫感
- 新增页面：AI 提供方管理、AI 聊天入口

### 5.8 配色方案

| 属性 | 浅色模式 | 暗色模式 |
|------|---------|---------|
| 页面背景 | `#fafafa` | `#1a1a2e` |
| 卡片背景 | `#ffffff` | `#252540` |
| 正文文字 | `#333333` | `#e0e0e0` |
| 辅助文字 | `#888888` | `#888888` |
| 标题文字 | `#1a1a1a` | `#f0f0f0` |
| 标签/强调色 | 多色点缀 | 多色点缀（提亮） |
| 链接色 | `#4078c0` | `#6db3f2` |
| 边框 | `#e8e8e8` | `#333350` |

### 5.9 字体与动效

- **字体栈**：`"PingFang SC", "Microsoft YaHei", "Helvetica Neue", Helvetica, Arial, sans-serif`
- **代码字体**：`"JetBrains Mono", "Fira Code", "Cascadia Code", Consolas, monospace`
- **动效**：
  - 卡片悬停：`transform: translateY(-2px)` + `box-shadow` 加深，过渡 200ms
  - 页面切换：淡入淡出 150ms
  - 搜索弹窗：缩放 + 透明度过渡 200ms
  - 时间轴节点：滚动渐入
- 不添加花哨的入场动画，保持克制

### 5.10 搜索

- 全屏覆盖弹窗
- 使用 Fuse.js 前端模糊搜索（文章标题 + 摘要）
- 无需额外后端接口（文章数据前端已有）
- 搜索结果实时展示，点击跳转

---

## 6. 实施顺序

| 阶段 | 内容 | 预估改动文件 |
|------|------|-------------|
| 1.1 | RSA+AES 加密通信 | 后端 5-8 个文件，前端 3-5 个文件 |
| 1.2 | 后台入口可配置 | 后端 2-3 个文件，前端 3-5 个文件 |
| 1.3 | 前端代码混淆 | 前端 1 个文件（vite.config.js） |
| 2.1 | 删除爬虫功能 | 删除 python-sdk/，后端 5-8 个文件 |
| 2.2 | AI 翻译功能 | 后端 8-12 个文件，前端 5-8 个文件 |
| 2.3 | AI 提供方管理 | 后端 5-8 个文件，前端 3-5 个文件 |
| 2.4 | AI 聊天框 | 后端 3-5 个文件，前端 3-5 个文件 |
| 3.1 | UI 全局布局 + 首页 | 前端 8-12 个文件 |
| 3.2 | 时间轴页面 | 后端 1-2 个文件，前端 3-5 个文件 |
| 3.3 | 文章详情/工具/后台 | 前端 5-8 个文件 |

---

## 7. 测试要点

### 安全层
- [ ] RSA 密钥对生成正确，公钥可被前端解析
- [ ] AES 密钥协商成功，Redis 存储正确
- [ ] 请求体加解密正确性（正常数据 + 边界数据）
- [ ] 密钥过期自动重新协商
- [ ] 白名单接口不受加密中间件影响
- [ ] 后台入口配置修改后路由正确响应
- [ ] 混淆后的前端代码可正常运行

### 功能层
- [ ] 爬虫相关代码完全移除，编译通过
- [ ] AI 翻译：发布时勾选 → 异步翻译 → 保存成功
- [ ] AI 翻译：编辑后 is_translated 重置
- [ ] AI 翻译：手动触发翻译/重新翻译
- [ ] 中英文切换：已翻译文章显示英文，未翻译显示中文+标注
- [ ] AI 提供方 CRUD + 连通性检测 + 余额查询
- [ ] AI 提供方 fallback 逻辑
- [ ] AI 聊天：多轮对话 + 流式返回

### UI 层
- [ ] 双栏布局在桌面/移动端表现正确
- [ ] 暗色模式切换全站同步
- [ ] 文章卡片流渲染正确
- [ ] 时间轴分组和排序正确
- [ ] 搜索功能正确
- [ ] 响应式适配

---

## 8. 新增依赖

### 后端（Go）
- 无需新增第三方依赖，全部使用标准库 + 现有依赖

### 前端（npm）
- `rollup-plugin-obfuscator` — 代码混淆
- `fuse.js` — 前端模糊搜索
