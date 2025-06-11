import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import Dashboard from './views/Dashboard';
import BlockExplorer from './views/BlockExplorer';
import Wallet from './views/Wallet';
import './assets/styles/App.css';

const { Header, Content, Footer } = Layout;

function App() {
  return (
    <Layout className="layout">
      <Header>
        <div className="logo" />
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
          ]}
        />
      </Header>
      <Content style={{ padding: '0 50px' }}>
        <div className="site-layout-content">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/explorer" element={<BlockExplorer />} />
            <Route path="/wallet" element={<Wallet />} />
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