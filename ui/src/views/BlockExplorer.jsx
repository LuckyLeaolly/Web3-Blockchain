import React, { useState, useEffect } from 'react';
import { Tabs, Input, Button, Table, Card, Typography, Collapse, Spin, message, Descriptions, Tag, Empty } from 'antd';
import { SearchOutlined, BlockOutlined, TransactionOutlined, WalletOutlined } from '@ant-design/icons';
import api from '../api/client';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Panel } = Collapse;

const BlockExplorer = () => {
  const [blocks, setBlocks] = useState([]);
  const [transactions, setTransactions] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState(null);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchPerformed, setSearchPerformed] = useState(false);
  const [loading, setLoading] = useState(true);

  // 加载初始数据
  useEffect(() => {
    const fetchInitialData = async () => {
      try {
        // 并行请求数据
        const [blocksResponse, transactionsResponse] = await Promise.all([
          api.blocks.getAll(10),
          api.transactions.getAll(10)
        ]);
        
        setBlocks(blocksResponse.data);
        setTransactions(transactionsResponse.data);
        setLoading(false);
      } catch (error) {
        console.error('获取数据失败:', error);
        message.error('加载区块链数据失败');
        setLoading(false);
      }
    };

    fetchInitialData();
  }, []);

  // 处理搜索
  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      message.warning('请输入搜索内容');
      return;
    }

    setSearchLoading(true);
    setSearchResults(null);
    setSearchPerformed(true);

    try {
      // 尝试按区块哈希搜索
      try {
        const blockResponse = await api.blocks.getById(searchQuery);
        setSearchResults({
          type: 'block',
          data: blockResponse.data
        });
        setSearchLoading(false);
        return;
      } catch (e) {
        // 不是区块，继续尝试其他类型
      }

      // 尝试按交易ID搜索
      try {
        const txResponse = await api.transactions.getById(searchQuery);
        setSearchResults({
          type: 'transaction',
          data: txResponse.data
        });
        setSearchLoading(false);
        return;
      } catch (e) {
        // 不是交易，继续尝试其他类型
      }

      // 尝试按地址搜索（余额）
      try {
        const balanceResponse = await api.wallets.getBalance(searchQuery);
        setSearchResults({
          type: 'address',
          data: {
            address: searchQuery,
            balance: balanceResponse.data.balance
          }
        });
        setSearchLoading(false);
        return;
      } catch (e) {
        // 不是有效地址
      }

      // 如果所有搜索都失败
      setSearchResults({ type: 'not_found' });
      message.info('未找到匹配结果');
    } catch (error) {
      console.error('搜索失败:', error);
      message.error('搜索过程中发生错误');
    } finally {
      setSearchLoading(false);
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
    if (!searchResults) return null;

    if (searchResults.type === 'block') {
      const block = searchResults.data;
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
    } else if (searchResults.type === 'transaction') {
      const tx = searchResults.data;
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
    } else if (searchResults.type === 'address') {
      return (
        <Card title="地址详情" style={{ marginTop: 16 }}>
          <p><strong>地址:</strong> {searchResults.data.address}</p>
          <p><strong>余额:</strong> {searchResults.data.balance}</p>
        </Card>
      );
    } else if (searchResults.type === 'not_found') {
      return (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="未找到匹配结果" />
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
              loading={searchLoading}
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