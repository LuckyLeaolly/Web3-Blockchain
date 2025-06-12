# API客户端使用指南

## 简介

本API客户端封装了与区块链后端API的所有交互，处理了认证、错误处理等常见问题，使前端组件可以更简洁地发起API请求。

## 主要功能

1. 自动添加认证Token
2. 处理认证失效自动跳转
3. 统一的错误处理
4. 按功能模块分类的API接口

## 使用方法

### 导入API客户端

```javascript
import api from '../api/client';
```

### 认证相关接口

```javascript
// 用户登录
const response = await api.auth.login({username: 'admin', password: 'admin123'});
const { token, userId } = response.data;

// 用户注册
await api.auth.register({username: 'newuser', email: 'user@example.com', password: 'password123'});
```

### 钱包相关接口

```javascript
// 获取所有钱包
const response = await api.wallets.getAll();
const addresses = response.data;

// 获取单个钱包余额
const balanceResponse = await api.wallets.getBalance('wallet-address');
const balance = balanceResponse.data.balance;

// 创建新钱包
const newWalletResponse = await api.wallets.create();
const newAddress = newWalletResponse.data.address;

// 获取钱包交易记录
const txResponse = await api.wallets.getTransactions('wallet-address');
const transactions = txResponse.data;
```

### 交易相关接口

```javascript
// 创建交易
const txData = {
  from: 'sender-address',
  to: 'receiver-address',
  amount: 10
};
const response = await api.transactions.create(txData);
const txId = response.data.txid;

// 获取所有交易
const allTx = await api.transactions.getAll(10); // 限制10条

// 获取特定交易
const tx = await api.transactions.getById('transaction-id');
```

### 区块相关接口

```javascript
// 获取所有区块
const blocks = await api.blocks.getAll(5); // 获取最新5个区块

// 获取特定区块
const block = await api.blocks.getById('block-hash');
```

### 系统信息接口

```javascript
// 获取区块链信息
const info = await api.system.getInfo();

// 获取网络配置
const config = await api.system.getNetworkConfig();

// 更新网络配置
await api.system.updateNetworkConfig({
  miningDifficulty: 4,
  miningReward: 50
});
```

## 错误处理

客户端会自动处理一些常见错误：

1. **401错误**：当认证失效时，会自动清除无效的token并跳转到登录页
2. **其他错误**：会被Promise.reject传递给调用者，可以在catch块中处理

```javascript
try {
  await api.transactions.create(txData);
} catch (error) {
  // 处理错误
  console.error('错误:', error);
  // 获取详细错误信息
  const errorMsg = error.response?.data?.error || '未知错误';
  message.error(errorMsg);
}
```

## 注意事项

1. API客户端会自动从localStorage读取认证token
2. 确保在发起需要认证的请求前，用户已经登录并获取了token
3. 使用try/catch块处理API请求可能出现的错误 