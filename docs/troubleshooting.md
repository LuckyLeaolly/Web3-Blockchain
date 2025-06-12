# Web3.0 区块链系统故障排除指南

## 接口错误问题

### 1. 交易创建接口返回500错误

当调用 `POST /api/v1/transactions` 接口创建交易时遇到500错误，可能有以下原因：

#### 问题原因：
1. **余额不足**：发送方钱包没有足够的余额进行转账
2. **交易验证失败**：交易签名验证失败或交易数据无效
3. **钱包不存在**：使用了不存在的钱包地址
4. **数据库锁定问题**：BadgerDB锁被其他进程占用

#### 解决方案：
1. 检查余额：确保发送方钱包有足够的余额
```bash
curl http://localhost:8080/api/v1/wallets/{wallet_address}/balance
```

2. 确认钱包存在：检查钱包地址是否存在
```bash
curl -H "Authorization: Bearer {token}" http://localhost:8080/api/v1/wallets
```

3. 重置数据库锁：如果是BadgerDB锁问题，可以尝试以下步骤：
```bash
# 停止所有相关进程
pkill -f "go run cmd/api/main.go"
# 删除锁文件
rm -f /path/to/web3-blockchain/data/blockchain/LOCK
# 重启服务
cd /path/to/web3-blockchain && go run cmd/api/main.go cmd/api/auth.go cmd/api/network.go
```

### 2. 区块链数据无法初始化

如果遇到"找不到区块链数据"或相关错误：

#### 解决方案：

1. 清空区块链数据并重新初始化：
```bash
rm -rf /path/to/web3-blockchain/data/blockchain/*
mkdir -p /path/to/web3-blockchain/data/blockchain
```

2. 修改代码，添加初始化逻辑，确保在区块链不存在时能够自动创建。

## 显示乱码问题

如果API返回的地址或其他数据显示为乱码：

#### 问题原因：
1. 公钥哈希直接输出导致的编码问题
2. 没有正确转换为Base58格式

#### 解决方案：
修改 `convertTransactionToOutput` 函数，确保地址使用Base58编码输出。

## BadgerDB相关问题

### 1. 无法获取数据库锁

```
Cannot acquire directory lock on ".../data/blockchain". Another process is using this Badger database.
```

#### 解决方案：

```bash
# 1. 查找并终止所有使用BadgerDB的进程
ps aux | grep blockchain
kill -9 [进程ID]

# 2. 删除锁文件
rm -f /path/to/web3-blockchain/data/blockchain/LOCK

# 3. 如果问题仍然存在，可以尝试重新启动机器
```

### 2. 数据损坏问题

如果遇到数据损坏警告或错误：

#### 解决方案：

```bash
# 备份现有数据（如果需要）
cp -r /path/to/web3-blockchain/data/blockchain /path/to/backup/

# 清空数据目录并重新初始化
rm -rf /path/to/web3-blockchain/data/blockchain/*
```

## 提高系统健壮性的建议

1. 添加适当的错误处理和恢复机制
2. 在API响应中提供详细的错误信息
3. 实现数据库自动恢复功能
4. 添加系统健康检查端点
5. 增强日志记录，包括详细的错误堆栈和上下文信息 