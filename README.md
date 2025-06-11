# Web3.0 区块链系统

这是一个使用 Golang 实现的 Web3.0 区块链系统，包括区块链核心和可视化的 Web 前端。

## 项目特性

- 使用 Go 语言开发的区块链核心
- 支持账户管理、交易处理、区块生成与同步
- 实现了基本的共识机制（PoW）
- 提供可视化的 Web 前端，包括区块链浏览器和简易钱包功能
- 实时展示区块链状态、交易信息
- 响应式设计，支持桌面和移动设备

## 技术栈

### 后端

- Golang
- Gin Web 框架
- BadgerDB (存储区块链数据)
- WebSocket (实时数据推送)

### 前端

- React
- Material UI / Ant Design
- Axios (HTTP 请求)
- Recharts (数据可视化)
- React Router (路由管理)
- ethers.js / web3.js (区块链交互)

## 目录结构

```
web3-blockchain/
├── cmd/                   # 应用程序入口点
│   ├── blockchain/        # 区块链节点入口
│   └── api/               # API 服务入口
├── internal/              # 内部包
│   ├── blockchain/        # 区块链核心实现
│   └── api/               # API 服务实现
├── pkg/                   # 公共包
│   ├── utils/             # 通用工具函数
│   └── models/            # 数据模型
├── ui/                    # 前端界面
│   ├── src/               # 源代码
│   │   ├── components/    # React组件
│   │   ├── pages/         # 页面组件
│   │   ├── hooks/         # 自定义React钩子
│   │   ├── context/       # React上下文
│   │   ├── utils/         # 工具函数
│   │   └── services/      # API服务调用
│   └── public/            # 静态资源
└── docs/                  # 文档
```

## 安装与运行

### 前提条件

- Go 1.16+
- Node.js 14+
- npm 或 yarn

### 安装步骤

1. 克隆仓库

```bash
git clone https://github.com/LuckyLeaolly/web3-blockchain.git
cd web3-blockchain
```

2. 安装后端依赖

```bash
go mod tidy
```

3. 安装前端依赖

```bash
cd ui
npm install # 或者使用 yarn install
```

4. 运行区块链节点

```bash
go run cmd/blockchain/main.go
```

5. 运行 API 服务

```bash
go run cmd/api/main.go
```

6. 运行前端界面

```bash
cd ui
npm start # 或者使用 yarn start
```

7. 打开浏览器访问 http://localhost:3000

## API 文档

API 文档使用 Swagger 生成，访问 http://localhost:8080/swagger/index.html

## 贡献

欢迎提交 Pull Request 或创建 Issue。

## 许可证

MIT 