# Blog 前端

基于 Vue 3 + Vite + Element Plus 的博客系统前端。

## 技术栈

- **框架**: Vue 3
- **构建工具**: Vite 5
- **UI组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP客户端**: Axios
- **CSS预处理器**: Less

## 目录结构

```
frontend/
├── src/
│   ├── assets/          # 静态资源
│   ├── components/      # 公共组件
│   ├── views/           # 页面视图
│   ├── router/          # 路由配置
│   ├── store/           # 状态管理
│   ├── api/             # API接口
│   ├── utils/           # 工具函数
│   └── main.js          # 入口文件
├── public/              # 公共资源
└── index.html           # HTML模板
```

## 快速开始

### 使用 Docker (推荐)

```bash
# 构建镜像
docker build -t blog-frontend .

# 运行容器
docker run -p 80:80 blog-frontend
```

### 本地开发

1. 安装依赖
```bash
npm install
```

2. 启动开发服务器
```bash
npm run dev
```

3. 构建生产版本
```bash
npm run build
```

## 环境变量

- `VITE_API_BASE_URL`: API基础URL

## 代码规范

```bash
# ESLint检查
npm run lint

# Prettier格式化
npm run format
```

## 许可证

MIT

