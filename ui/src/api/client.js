import axios from 'axios';

// API基础URL
const API_URL = 'http://localhost:8080/api/v1';

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_URL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  }
});

// 请求拦截器，添加认证token
apiClient.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    console.error('请求错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器，处理常见错误
apiClient.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    // 处理401未授权错误
    if (error.response && error.response.status === 401) {
      // 清除无效的token
      localStorage.removeItem('token');
      localStorage.removeItem('userId');
      
      // 重定向到登录页面
      window.location.href = '/login';
      return Promise.reject(new Error('身份验证失败，请重新登录'));
    }
    
    return Promise.reject(error);
  }
);

// API方法
const api = {
  // 认证相关
  auth: {
    login: (credentials) => apiClient.post('/auth/login', credentials),
    register: (userData) => apiClient.post('/auth/register', userData),
  },
  
  // 钱包相关
  wallets: {
    getAll: () => apiClient.get('/wallets'),
    getBalance: (address) => apiClient.get(`/wallets/${address}/balance`),
    create: () => apiClient.post('/wallets'),
    getTransactions: (address) => apiClient.get(`/wallets/${address}/transactions`),
  },
  
  // 交易相关
  transactions: {
    create: (txData) => apiClient.post('/transactions', txData),
    getAll: (limit = 10) => apiClient.get(`/transactions?limit=${limit}`),
    getById: (txId) => apiClient.get(`/transactions/${txId}`),
  },
  
  // 区块相关
  blocks: {
    getAll: (limit = 10) => apiClient.get(`/blocks?limit=${limit}`),
    getById: (blockId) => apiClient.get(`/blocks/${blockId}`),
  },
  
  // 系统信息
  system: {
    getInfo: () => apiClient.get('/info'),
    getNetworkConfig: () => apiClient.get('/network/config'),
    updateNetworkConfig: (config) => apiClient.put('/network/config', config),
  }
};

export { API_URL };
export default api; 