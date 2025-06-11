# Web3.0 区块链系统

一个基于Go语言开发的基础区块链系统，带有Web前端界面。本项目实现了区块链的核心功能，包括区块生成、工作量证明、交易处理和钱包管理。同时提供了RESTful API和前端界面，方便用户交互和数据监控。

## 项目结构

```
web3-blockchain/
├── cmd/                  # 命令行入口
│   ├── blockchain/       # 区块链节点CLI
│   └── api/              # API服务
├── pkg/                  # 核心包
│   └── models/           # 数据模型
│       ├── block.go      # 区块结构
│       ├── blockchain.go # 区块链
│       ├── proof_of_work.go # 工作量证明
│       ├── transaction.go # 交易
│       └── wallet.go     # 钱包
├── ui/                   # 前端界面
│   └── src/              # React源代码
└── docs/                 # API文档
    └── swagger.json      # Swagger文档
```

## 功能特性

- **区块链核心功能**
  - 基于工作量证明(PoW)的共识机制
  - UTXO交易模型
  - 多钱包管理
  - 数据持久化存储

- **API服务**
  - RESTful API接口
  - Swagger API文档
  - 区块链浏览功能
  - 交易创建和查询

- **Web前端**
  - 区块链数据可视化
  - 区块和交易浏览器
  - 钱包管理界面
  - 发送交易功能

## 环境要求

- Go 1.18+
- Node.js 14+
- BadgerDB (自动安装)

## 安装步骤

### 1. 克隆仓库

```bash
git clone https://github.com/yourname/web3-blockchain.git
cd web3-blockchain
```

### 2. 安装Go依赖

```bash
go mod tidy
```

### 3. 安装前端依赖

```bash
cd ui
npm install
```

## 使用说明

### 运行区块链节点

```bash
cd cmd/blockchain
go run main.go
```

常用命令:
- `go run main.go createblockchain -address ADDRESS` - 创建一个新的区块链
- `go run main.go createwallet` - 创建一个新的钱包
- `go run main.go getbalance -address ADDRESS` - 获取地址余额
- `go run main.go send -from FROM -to TO -amount AMOUNT` - 发送代币
- `go run main.go printchain` - 打印区块链内容

### 运行API服务

```bash
cd cmd/api
go run main.go
```

API服务将在 http://localhost:8080 上启动，Swagger文档可通过 http://localhost:8080/swagger/index.html 访问。

### 运行前端界面

```bash
cd ui
npm start
```

前端界面将在 http://localhost:3000 上启动。

## 技术栈

- **后端**: Go, BadgerDB
- **API框架**: Gin
- **文档**: Swagger
- **前端**: React, Ant Design

## 未来计划

- 实现更高效的共识机制(PoS)
- 添加智能合约功能
- 扩展P2P网络功能
- 优化区块同步机制

## 贡献指南

1. Fork项目
2. 创建特性分支: `git checkout -b feature/new-feature`
3. 提交更改: `git commit -am 'Add new feature'`
4. 推送到分支: `git push origin feature/new-feature`
5. 提交Pull Request

## 许可证

MIT 