# Web3.0 区块链系统

一个基于Go语言开发的基础区块链系统，带有Web前端界面。本项目实现了区块链的核心功能，包括区块生成、工作量证明、交易处理和钱包管理。同时提供了RESTful API和前端界面，方便用户交互和数据监控。

## 重要说明：数据文件管理

**注意：** 本项目使用BadgerDB存储区块链数据，会生成较大的数据文件。这些文件已被添加到`.gitignore`中，不应该提交到版本控制系统。

- 数据文件位置：`data/blockchain/`目录
- 详细说明：查看[数据管理指南](docs/data-management.md)了解更多信息

如果你不小心将这些大文件推送到了GitHub，请使用项目根目录下的`clean_git.sh`脚本清理Git历史：

```bash
chmod +x clean_git.sh
./clean_git.sh
```

## 项目结构

```
web3-blockchain/
├── cmd/                  # 命令行入口
│   ├── blockchain/       # 区块链节点CLI
│   └── api/              # API服务
│       ├── main.go       # 主程序入口
│       ├── auth.go       # 认证相关接口
│       └── network.go    # 网络配置接口
├── pkg/                  # 核心包
│   ├── common/           # 通用工具
│   │   └── paths.go      # 路径管理
│   └── models/           # 数据模型
│       ├── block.go      # 区块结构
│       ├── blockchain.go # 区块链
│       ├── proof_of_work.go # 工作量证明
│       ├── transaction.go # 交易
│       └── wallet.go     # 钱包
├── ui/                   # 前端界面
│   └── src/              # React源代码
├── data/                 # 数据存储目录 
│   ├── blockchain/       # 区块链数据
│   └── wallets/          # 钱包数据
└── docs/                 # API文档
    └── swagger.json      # Swagger文档
```

## 功能特性

- **区块链核心功能**
  - 基于工作量证明(PoW)的共识机制
  - UTXO交易模型
  - 多钱包管理
  - 安全的钱包序列化与存储
  - 数据持久化存储 (BadgerDB)
  - 安全的交易签名与验证
  - 统一的数据路径管理

- **认证与安全**
  - JWT认证系统
  - API访问控制
  - 用户登录/注册功能
  - 受保护资源访问

- **API服务**
  - 完整的RESTful API接口
  - 详细的Swagger API文档
  - 区块链信息查询（区块高度、总交易数、网络状态）
  - 区块相关接口（获取区块列表、按哈希或高度查询区块）
  - 交易相关接口（创建交易、查询交易详情、获取交易列表）
  - 钱包管理接口（创建钱包、获取余额、查询钱包地址列表）
  - 网络配置接口（获取和更新共识参数）

- **Web前端**
  - **仪表盘**：展示区块链关键指标（区块高度、交易总数、运行状态）
  - **区块浏览器**：
    - 区块和交易列表展示
    - 强大的搜索功能（支持区块哈希、交易ID和钱包地址搜索）
    - 区块与交易的详细信息查看
  - **钱包管理**：
    - 创建新钱包
    - 查看钱包余额
    - 发送交易功能
  - **交易历史**：
    - 按钱包地址查询交易记录
    - 按日期范围筛选交易
    - 交易时间线可视化展示
    - 收入/支出交易区分
  - **网络配置**：
    - 查看和修改网络参数
    - 区块生成时间调整
    - 难度调整
    - 挖矿奖励配置

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

### 1. 创建钱包

首先，创建一个钱包地址：

```bash
cd cmd/blockchain
go run main.go createwallet
```

系统将生成一个新的钱包地址，记录该地址以便后续使用。

### 2. 初始化区块链

使用刚才创建的钱包地址初始化区块链：

```bash
go run main.go createblockchain -address YOUR_WALLET_ADDRESS
```

此命令将创建一个新的区块链数据库，并生成创世区块，该地址将获得初始挖矿奖励。

### 3. 查询钱包余额

```bash
go run main.go getbalance -address YOUR_WALLET_ADDRESS
```

### 4. 运行API服务

```bash
cd cmd/api
go run main.go auth.go network.go
```

API服务将在 http://localhost:8080 上启动，Swagger文档可通过 http://localhost:8080/swagger/index.html 访问。

### 5. 启动前端界面

```bash
cd ui
npm start
```

前端界面将在 http://localhost:3000 上启动。

### 6. API认证

要访问受保护的API端点，需要先登录获取JWT令牌：

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"admin", "password":"admin123"}' \
  http://localhost:8080/api/v1/auth/login
```

使用返回的令牌访问受保护的API：

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/network/config
```

## 技术栈

- **后端**: Go, BadgerDB
- **API框架**: Gin
- **认证**: JWT (JSON Web Tokens)
- **数据库**: BadgerDB (键值存储)
- **文档**: Swagger
- **前端**: React, Ant Design
- **状态管理**: React Hooks
- **API通信**: Axios
- **时间处理**: Moment.js

## 系统组件详解

### 区块链核心

- **区块结构**：包含前一区块哈希、时间戳、随机数、交易列表等
- **工作量证明**：确保区块的生成需要一定的计算工作，防止恶意攻击
- **钱包管理**：生成ECDSA密钥对，创建基于公钥的地址，安全存储
- **UTXO模型**：使用未花费交易输出模型处理交易，防止双花

### API服务

- **RESTful设计**：遵循REST原则设计的API接口
- **认证中间件**：基于JWT的API访问控制
- **路由管理**：清晰的路由结构，区分公开和需要认证的API

### 数据存储

- **统一路径管理**：通过common包集中管理所有数据存储路径
- **BadgerDB存储**：高性能的键值对存储，适合区块链数据
- **序列化优化**：改进对象序列化方法，确保数据一致性

## 未来计划

- 实现更高效的共识机制(PoS)
- 添加智能合约功能
- 扩展P2P网络功能
- 优化区块同步机制
- 添加更丰富的数据可视化图表
- 实现HD钱包（分层确定性钱包）功能
- 增强安全性和隐私保护

## 问题排查与解决

如果遇到常见问题，请参考以下解决方案：

1. **钱包创建问题**：如遇"gob: type elliptic.p256Curve has no exported fields"错误，已通过改进序列化方法解决
2. **数据路径问题**：所有数据路径已统一管理，确保跨平台兼容性
3. **API认证失败**：检查JWT令牌格式和有效期
4. **区块链访问错误**：确保BadgerDB数据库文件未损坏，必要时重新初始化区块链

## 贡献指南

1. Fork项目
2. 创建特性分支: `git checkout -b feature/new-feature`
3. 提交更改: `git commit -am 'Add new feature'`
4. 推送到分支: `git push origin feature/new-feature`
5. 提交Pull Request

## 许可证

MIT 