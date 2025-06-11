import React, { useState, useEffect } from 'react';
import { Table, Card, Typography, Spin, message, Timeline, Tag, Select, DatePicker, Button } from 'antd';
import { ClockCircleOutlined, ArrowUpOutlined, ArrowDownOutlined, SearchOutlined } from '@ant-design/icons';
import axios from 'axios';
import moment from 'moment';

const { Title, Text } = Typography;
const { Option } = Select;
const { RangePicker } = DatePicker;

const TransactionHistory = () => {
  const [loading, setLoading] = useState(true);
  const [transactions, setTransactions] = useState([]);
  const [walletAddress, setWalletAddress] = useState('');
  const [walletAddresses, setWalletAddresses] = useState([]);
  const [dateRange, setDateRange] = useState(null);
  const [filteredTransactions, setFilteredTransactions] = useState([]);

  const API_URL = 'http://localhost:8080/api/v1';

  // 获取钱包地址列表
  const fetchWalletAddresses = async () => {
    try {
      const response = await axios.get(`${API_URL}/wallets`);
      setWalletAddresses(response.data);
      if (response.data.length > 0) {
        setWalletAddress(response.data[0]);
      }
    } catch (error) {
      console.error('获取钱包列表失败:', error);
      message.error('获取钱包列表失败');
    }
  };

  // 获取交易历史
  const fetchTransactions = async (address) => {
    if (!address) return;

    setLoading(true);
    try {
      const response = await axios.get(`${API_URL}/wallets/${address}/transactions`);
      const sortedTransactions = response.data.sort((a, b) => b.timestamp - a.timestamp);
      setTransactions(sortedTransactions);
      setFilteredTransactions(sortedTransactions);
      setLoading(false);
    } catch (error) {
      console.error('获取交易历史失败:', error);
      message.error('获取交易历史失败');
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchWalletAddresses();
  }, []);

  useEffect(() => {
    if (walletAddress) {
      fetchTransactions(walletAddress);
    }
  }, [walletAddress]);

  // 处理日期范围变化
  const handleDateRangeChange = (dates) => {
    setDateRange(dates);
  };

  // 过滤交易记录
  const filterTransactions = () => {
    let filtered = [...transactions];

    // 根据日期范围过滤
    if (dateRange && dateRange[0] && dateRange[1]) {
      const startTime = dateRange[0].startOf('day').valueOf() / 1000;
      const endTime = dateRange[1].endOf('day').valueOf() / 1000;

      filtered = filtered.filter(tx => 
        tx.timestamp >= startTime && tx.timestamp <= endTime
      );
    }

    setFilteredTransactions(filtered);
  };

  // 交易列表列定义
  const columns = [
    {
      title: '交易ID',
      dataIndex: 'id',
      key: 'id',
      render: id => (
        <Text ellipsis={{ tooltip: id }}>
          {id.substring(0, 12)}...
        </Text>
      ),
    },
    {
      title: '类型',
      key: 'type',
      render: (_, record) => {
        const isSender = record.from === walletAddress;
        return (
          <Tag color={isSender ? "volcano" : "green"}>
            {isSender ? '转出' : '转入'}
          </Tag>
        );
      },
    },
    {
      title: '对方地址',
      key: 'counterparty',
      render: (_, record) => {
        const counterparty = record.from === walletAddress ? record.to : record.from;
        return (
          <Text ellipsis={{ tooltip: counterparty }}>
            {counterparty.substring(0, 16)}...
          </Text>
        );
      },
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount, record) => {
        const isSender = record.from === walletAddress;
        const style = { color: isSender ? '#cf1322' : '#3f8600' };
        return (
          <Text style={style}>
            {isSender ? '-' : '+'}{amount}
          </Text>
        );
      },
    },
    {
      title: '时间',
      dataIndex: 'timestamp',
      key: 'timestamp',
      render: timestamp => moment(timestamp * 1000).format('YYYY-MM-DD HH:mm:ss'),
    },
  ];

  if (loading && !walletAddress) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
        <p>加载钱包数据中...</p>
      </div>
    );
  }

  return (
    <div>
      <Title level={2}>交易历史</Title>
      
      <Card style={{ marginBottom: 16 }}>
        <div style={{ marginBottom: 16 }}>
          <Text strong>选择钱包地址: </Text>
          <Select 
            style={{ width: 360, marginRight: 16 }}
            value={walletAddress}
            onChange={setWalletAddress}
          >
            {walletAddresses.map(address => (
              <Option key={address} value={address}>{address}</Option>
            ))}
          </Select>
          
          <Text strong style={{ marginLeft: 16, marginRight: 8 }}>日期范围: </Text>
          <RangePicker onChange={handleDateRangeChange} style={{ marginRight: 16 }} />
          
          <Button 
            type="primary" 
            icon={<SearchOutlined />}
            onClick={filterTransactions}
          >
            筛选
          </Button>
        </div>
        
        {loading ? (
          <div style={{ textAlign: 'center', padding: '20px' }}>
            <Spin />
            <p>加载交易历史中...</p>
          </div>
        ) : filteredTransactions.length > 0 ? (
          <>
            <Table 
              dataSource={filteredTransactions} 
              columns={columns} 
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
            
            <Card title="交易时间线" style={{ marginTop: 16 }}>
              <Timeline mode="left">
                {filteredTransactions.slice(0, 10).map(tx => {
                  const isSender = tx.from === walletAddress;
                  return (
                    <Timeline.Item
                      key={tx.id}
                      color={isSender ? 'red' : 'green'}
                      dot={isSender ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
                      label={moment(tx.timestamp * 1000).format('YYYY-MM-DD HH:mm:ss')}
                    >
                      <p>
                        <Tag color={isSender ? "volcano" : "green"}>
                          {isSender ? '转出' : '转入'}
                        </Tag>
                        {isSender ? `转给 ${tx.to.substring(0, 8)}...` : `来自 ${tx.from.substring(0, 8)}...`}
                      </p>
                      <p>金额: <Text strong>{tx.amount}</Text></p>
                      <p>交易ID: <Text type="secondary">{tx.id.substring(0, 16)}...</Text></p>
                    </Timeline.Item>
                  );
                })}
              </Timeline>
            </Card>
          </>
        ) : (
          <div style={{ textAlign: 'center', padding: '20px' }}>
            <Text type="secondary">暂无交易记录</Text>
          </div>
        )}
      </Card>
    </div>
  );
};

export default TransactionHistory; 