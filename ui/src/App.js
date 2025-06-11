import React, { useState, useEffect } from 'react';
import { Routes, Route, Link, Navigate, useNavigate } from 'react-router-dom';
import { Layout, Menu, Button, Avatar, Dropdown } from 'antd';
import { UserOutlined, SettingOutlined, LogoutOutlined } from '@ant-design/icons';
import Dashboard from './views/Dashboard';
import BlockExplorer from './views/BlockExplorer';
import Wallet from './views/Wallet';
import Login from './views/Login';
import NetworkConfig from './views/NetworkConfig';
import TransactionHistory from './views/TransactionHistory';
import './assets/styles/App.css';

const { Header, Content, Footer } = Layout;

// 受保护的路由组件
const ProtectedRoute = ({ children, isAuthenticated }) => {
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return children;
};

function App() {
  const [auth, setAuth] = useState({
    isAuthenticated: false,
    token: null,
    userId: null,
  });
  
  const navigate = useNavigate();

  // 初始化时检查本地存储中是否有令牌
  useEffect(() => {
    const token = localStorage.getItem('token');
    const userId = localStorage.getItem('userId');
    
    if (token && userId) {
      setAuth({
        isAuthenticated: true,
        token,
        userId,
      });
    }
  }, []);

  // 处理注销
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    setAuth({
      isAuthenticated: false,
      token: null,
      userId: null,
    });
    navigate('/login');
  };

  // 用户菜单项
  const userMenuItems = [
    {
      key: '1',
      label: '账户设置',
      icon: <SettingOutlined />,
    },
    {
      key: '2',
      label: '退出登录',
      icon: <LogoutOutlined />,
      onClick: handleLogout,
    },
  ];

  return (
    <Layout className="layout">
      <Header>
        <div className="logo" />
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
          <Menu
            theme="dark"
            mode="horizontal"
            defaultSelectedKeys={['1']}
            items={[
              {
                key: '1',
                label: <Link to="/">首页</Link>,
              },
              {
                key: '2',
                label: <Link to="/explorer">区块浏览器</Link>,
              },
              {
                key: '3',
                label: <Link to="/wallet">钱包</Link>,
              },
              {
                key: '4',
                label: <Link to="/transactions">交易历史</Link>,
              },
              {
                key: '5',
                label: <Link to="/network">网络配置</Link>,
              },
            ]}
          />
          
          {auth.isAuthenticated ? (
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Avatar icon={<UserOutlined />} style={{ cursor: 'pointer', backgroundColor: '#1890ff' }} />
            </Dropdown>
          ) : (
            <Button type="primary" onClick={() => navigate('/login')}>登录</Button>
          )}
        </div>
      </Header>
      <Content style={{ padding: '0 50px' }}>
        <div className="site-layout-content">
          <Routes>
            <Route path="/login" element={<Login setAuth={setAuth} />} />
            <Route path="/" element={
              <ProtectedRoute isAuthenticated={auth.isAuthenticated}>
                <Dashboard />
              </ProtectedRoute>
            } />
            <Route path="/explorer" element={
              <ProtectedRoute isAuthenticated={auth.isAuthenticated}>
                <BlockExplorer />
              </ProtectedRoute>
            } />
            <Route path="/wallet" element={
              <ProtectedRoute isAuthenticated={auth.isAuthenticated}>
                <Wallet />
              </ProtectedRoute>
            } />
            <Route path="/transactions" element={
              <ProtectedRoute isAuthenticated={auth.isAuthenticated}>
                <TransactionHistory />
              </ProtectedRoute>
            } />
            <Route path="/network" element={
              <ProtectedRoute isAuthenticated={auth.isAuthenticated}>
                <NetworkConfig />
              </ProtectedRoute>
            } />
          </Routes>
        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        Web3.0 区块链系统 ©{new Date().getFullYear()}
      </Footer>
    </Layout>
  );
}

export default App; 