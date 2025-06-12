import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic, Table, Typography, Spin } from 'antd';
import { BlockOutlined, TransactionOutlined, NodeIndexOutlined } from '@ant-design/icons';
import api from '../api/client';
import { Link } from 'react-router-dom';

const { Title } = Typography;

const Dashboard = () => {
  const [loading, setLoading] = useState(true);
  const [blockchainInfo, setBlockchainInfo] = useState({
    height: 0,
    transactions: 0,
    status: 'loading',
    version: '0.0.0',
  });
  const [latestBlocks, setLatestBlocks] = useState([]);
  const [latestTransactions, setLatestTransactions] = useState([]);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        // 并行请求数据
        const [infoResponse, blocksResponse, transactionsResponse] = await Promise.all([
          api.system.getInfo(),
          api.blocks.getAll(5), // 获取最新的5个区块
          api.transactions.getAll(5), // 获取最新的5个交易
        ]);
        
        setBlockchainInfo(infoResponse.data);
        setLatestBlocks(blocksResponse.data);
        setLatestTransactions(transactionsResponse.data);
        setLoading(false);
      } catch (error) {
        console.error('获取仪表盘数据失败:', error);
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  const blockColumns = [
    {
      title: '高度',
      dataIndex: 'height',
      key: 'height',
    },
    {
      title: '哈希',
      dataIndex: 'hash',
      key: 'hash',
      render: (hash) => `${hash.substring(0, 16)}...`,
    },
    {
      title: '交易数',
      dataIndex: 'transactions',
      key: 'transactions',
      render: (txs) => txs.length,
    },
    {
      title: '时间戳',
      dataIndex: 'timestamp',
      key: 'timestamp',
      render: (timestamp) => new Date(timestamp * 1000).toLocaleString(),
    },
  ];

  const txColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      render: (id) => `${id.substring(0, 16)}...`,
    },
    {
      title: '发送方',
      dataIndex: 'from',
      key: 'from',
      render: (from) => (from === '系统' ? from : `${from.substring(0, 12)}...`),
    },
    {
      title: '接收方',
      dataIndex: 'to',
      key: 'to',
      render: (to) => `${to.substring(0, 12)}...`,
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
    },
  ];

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
        <p>加载仪表盘数据中...</p>
      </div>
    );
  }

  return (
    <div>
      <Title level={2}>区块链状态</Title>
      <div className="dashboard-stats">
        <Row gutter={16}>
          <Col span={8}>
            <Card>
              <Statistic
                title="当前区块高度"
                value={blockchainInfo.height}
                prefix={<BlockOutlined />}
              />
            </Card>
          </Col>
          <Col span={8}>
            <Card>
              <Statistic
                title="交易总数"
                value={blockchainInfo.transactions}
                prefix={<TransactionOutlined />}
              />
            </Card>
          </Col>
          <Col span={8}>
            <Card>
              <Statistic
                title="运行状态"
                value={blockchainInfo.status}
                prefix={<NodeIndexOutlined />}
              />
            </Card>
          </Col>
        </Row>
      </div>

      <div style={{ marginTop: 24 }}>
        <Title level={3}>最新区块</Title>
        <Table 
          dataSource={latestBlocks} 
          columns={blockColumns}
          rowKey="hash"
          pagination={false}
        />
      </div>

      <div style={{ marginTop: 24 }}>
        <Title level={3}>最新交易</Title>
        <Table 
          dataSource={latestTransactions} 
          columns={txColumns}
          rowKey="id"
          pagination={false}
        />
      </div>
    </div>
  );
};

export default Dashboard; 