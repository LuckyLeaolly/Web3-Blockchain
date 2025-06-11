import React, { useState, useEffect } from 'react';
import { Tabs, Input, Button, Table, Card, Typography, Collapse, Spin, message } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import axios from 'axios';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Panel } = Collapse;

const BlockExplorer = () => {
  const [blocks, setBlocks] = useState([]);
  const [transactions, setTransactions] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResult, setSearchResult] = useState(null);
  const [searching, setSearching] = useState(false);
  const [loading, setLoading] = useState(true);

  const API_URL = 'http://localhost:8080/api/v1';

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        
        // 获取区块列表
        const blocksResponse = await axios.get(`${API_URL}/blocks?limit=10`);
        setBlocks(blocksResponse.data);
        
        // 获取交易列表
        const transactionsResponse = await axios.get(`${API_URL}/transactions?limit=10`);
        setTransactions(transactionsResponse.data);
        
        setLoading(false);
      } catch (error) {
        console.error('获取数据失败:', error);
        setLoading(false);
        message.error('加载数据失败，请稍后再试');
      }
    };

    fetchData();
  }, []);

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      message.warning('请输入搜索内容');
      return;
    }

    setSearching(true);
    setSearchResult(null);

    try {
      // 尝试按区块哈希查询
      if (searchQuery.length >= 32) {
        try {
          const blockResponse = await axios.get(`${API_URL}/blocks/${searchQuery}`);
          setSearchResult({ type: 'block', data: blockResponse.data });
          setSearching(false);
          return;
        } catch (error) {
          // 如果不是区块，可能是交易
          if (error.response && error.response.status !== 404) {
            throw error;
          }
        }
      }

      // 尝试按交易ID查询
      try {
        const txResponse = await axios.get(`${API_URL}/transactions/${searchQuery}`);
        setSearchResult({ type: 'transaction', data: txResponse.data });
        setSearching(false);
        return;
      } catch (error) {
        // 如果不是交易，可能是其他内容
        if (error.response && error.response.status !== 404) {
          throw error;
        }
      }

      // 尝试按地址查询余额
      try {
        const balanceResponse = await axios.get(`${API_URL}/wallets/${searchQuery}/balance`);
        setSearchResult({ type: 'address', data: { address: searchQuery, ...balanceResponse.data } });
        setSearching(false);
        return;
      } catch (error) {
        if (error.response && error.response.status !== 404) {
          throw error;
        }
      }

      // 如果都未找到
      message.warning('未找到匹配的区块、交易或地址');
      setSearching(false);
    } catch (error) {
      console.error('搜索失败:', error);
      message.error('搜索失败，请稍后再试');
      setSearching(false);
    }
  };

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
      render: hash => (
        <Text ellipsis={{ tooltip: hash }}>
          {hash}
        </Text>
      ),
    },
    {
      title: '时间戳',
      dataIndex: 'timestamp',
      key: 'timestamp',
      render: timestamp => new Date(timestamp * 1000).toLocaleString(),
    },
    {
      title: '交易数量',
      dataIndex: 'transactions',
      key: 'transactions',
      render: txs => txs.length,
    },
  ];

  const transactionColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      render: id => (
        <Text ellipsis={{ tooltip: id }}>
          {id}
        </Text>
      ),
    },
    {
      title: '发送方',
      dataIndex: 'from',
      key: 'from',
      render: from => (
        <Text ellipsis={{ tooltip: from }}>
          {from}
        </Text>
      ),
    },
    {
      title: '接收方',
      dataIndex: 'to',
      key: 'to',
      render: to => (
        <Text ellipsis={{ tooltip: to }}>
          {to}
        </Text>
      ),
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
    },
  ];

  // 渲染搜索结果
  const renderSearchResult = () => {
    if (!searchResult) return null;

    if (searchResult.type === 'block') {
      const block = searchResult.data;
      return (
        <Card title={`区块详情 (高度: ${block.height})`} style={{ marginTop: 16 }}>
          <p><strong>哈希:</strong> {block.hash}</p>
          <p><strong>前一个区块哈希:</strong> {block.prevBlockHash}</p>
          <p><strong>时间戳:</strong> {new Date(block.timestamp * 1000).toLocaleString()}</p>
          <p><strong>Nonce:</strong> {block.nonce}</p>
          
          <Collapse>
            <Panel header={`交易 (${block.transactions.length})`} key="1">
              <Table 
                dataSource={block.transactions} 
                columns={transactionColumns}
                rowKey="id"
              />
            </Panel>
          </Collapse>
        </Card>
      );
    } else if (searchResult.type === 'transaction') {
      const tx = searchResult.data;
      return (
        <Card title="交易详情" style={{ marginTop: 16 }}>
          <p><strong>交易ID:</strong> {tx.id}</p>
          <p><strong>发送方:</strong> {tx.from}</p>
          <p><strong>接收方:</strong> {tx.to}</p>
          <p><strong>金额:</strong> {tx.amount}</p>
          <p><strong>时间戳:</strong> {new Date(tx.timestamp * 1000).toLocaleString()}</p>
          
          <Collapse>
            <Panel header="输入" key="1">
              <ul>
                {tx.inputs.map((input, index) => (
                  <li key={index}>{input}</li>
                ))}
              </ul>
            </Panel>
          </Collapse>
        </Card>
      );
    } else if (searchResult.type === 'address') {
      return (
        <Card title="地址详情" style={{ marginTop: 16 }}>
          <p><strong>地址:</strong> {searchResult.data.address}</p>
          <p><strong>余额:</strong> {searchResult.data.balance}</p>
        </Card>
      );
    }
    
    return null;
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
        <p>加载数据中...</p>
      </div>
    );
  }

  return (
    <div>
      <Title level={2}>区块浏览器</Title>
      
      <div className="search-section">
        <Input.Search
          placeholder="输入区块哈希、交易ID或地址进行搜索"
          enterButton={
            <Button 
              icon={<SearchOutlined />} 
              loading={searching}
              type="primary"
            >
              搜索
            </Button>
          }
          size="large"
          value={searchQuery}
          onChange={e => setSearchQuery(e.target.value)}
          onSearch={handleSearch}
          style={{ marginBottom: 16 }}
        />
        
        {renderSearchResult()}
      </div>
      
      <Tabs defaultActiveKey="1" style={{ marginTop: 24 }}>
        <TabPane tab="区块列表" key="1">
          <Table 
            dataSource={blocks} 
            columns={blockColumns}
            rowKey="hash"
          />
        </TabPane>
        <TabPane tab="交易列表" key="2">
          <Table 
            dataSource={transactions} 
            columns={transactionColumns}
            rowKey="id"
          />
        </TabPane>
      </Tabs>
    </div>
  );
};

export default BlockExplorer; 