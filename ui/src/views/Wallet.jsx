import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, Select, Typography, Table, message, Tabs, Spin } from 'antd';
import { WalletOutlined, SendOutlined, PlusOutlined } from '@ant-design/icons';
import axios from 'axios';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;

const Wallet = () => {
  const [addresses, setAddresses] = useState([]);
  const [selectedAddress, setSelectedAddress] = useState('');
  const [balance, setBalance] = useState(0);
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [sending, setSending] = useState(false);
  const [transactions, setTransactions] = useState([]);

  const API_URL = 'http://localhost:8080/api/v1';

  // 加载钱包地址列表
  const fetchAddresses = async () => {
    try {
      const response = await axios.get(`${API_URL}/wallets`);
      setAddresses(response.data);
      if (response.data.length > 0 && !selectedAddress) {
        setSelectedAddress(response.data[0]);
        fetchBalance(response.data[0]);
      }
      setLoading(false);
    } catch (error) {
      console.error('获取钱包列表失败:', error);
      message.error('获取钱包列表失败');
      setLoading(false);
    }
  };

  // 获取余额
  const fetchBalance = async (address) => {
    if (!address) return;

    try {
      const response = await axios.get(`${API_URL}/wallets/${address}/balance`);
      setBalance(response.data.balance);
    } catch (error) {
      console.error('获取余额失败:', error);
      message.error('获取余额失败');
    }
  };

  useEffect(() => {
    fetchAddresses();
  }, []);

  // 当选择的地址变化时，获取新的余额
  useEffect(() => {
    if (selectedAddress) {
      fetchBalance(selectedAddress);
    }
  }, [selectedAddress]);

  // 创建新钱包
  const handleCreateWallet = async () => {
    setCreating(true);
    try {
      const response = await axios.post(`${API_URL}/wallets`);
      const newAddress = response.data.address;
      
      message.success(`新钱包创建成功: ${newAddress}`);
      setAddresses([...addresses, newAddress]);
      setSelectedAddress(newAddress);
      fetchBalance(newAddress);
      setCreating(false);
    } catch (error) {
      console.error('创建钱包失败:', error);
      message.error('创建钱包失败');
      setCreating(false);
    }
  };

  // 发送交易
  const handleSend = async (values) => {
    setSending(true);
    try {
      const txData = {
        from: selectedAddress,
        to: values.toAddress,
        amount: parseInt(values.amount),
      };
      
      const response = await axios.post(`${API_URL}/transactions`, txData);
      
      message.success(`交易已创建: ${response.data.txid}`);
      form.resetFields(['toAddress', 'amount']);
      
      // 重新获取余额和交易记录
      fetchBalance(selectedAddress);
      setSending(false);
    } catch (error) {
      console.error('发送交易失败:', error);
      message.error('发送交易失败: ' + (error.response?.data?.error || '未知错误'));
      setSending(false);
    }
  };

  // 当选择的地址发生变化
  const handleAddressChange = (value) => {
    setSelectedAddress(value);
    form.setFieldsValue({ fromAddress: value });
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
        <p>加载钱包数据中...</p>
      </div>
    );
  }

  return (
    <div className="wallet-container">
      <Title level={2}>区块链钱包</Title>
      
      <Tabs defaultActiveKey="1">
        <TabPane 
          tab={
            <span>
              <WalletOutlined />
              我的钱包
            </span>
          } 
          key="1"
        >
          <Card title="钱包管理" style={{ marginBottom: 16 }}>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text strong>选择地址:</Text> 
                <Select 
                  style={{ width: 360, marginLeft: 8 }}
                  value={selectedAddress}
                  onChange={handleAddressChange}
                >
                  {addresses.map(address => (
                    <Option key={address} value={address}>{address}</Option>
                  ))}
                </Select>
              </div>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={handleCreateWallet}
                loading={creating}
              >
                创建新钱包
              </Button>
            </div>
            
            {selectedAddress && (
              <div style={{ marginBottom: 16 }}>
                <Card type="inner" title="钱包详情">
                  <p><strong>地址:</strong> {selectedAddress}</p>
                  <p><strong>余额:</strong> {balance}</p>
                </Card>
              </div>
            )}
          </Card>
          
          <Card title="发送交易" style={{ marginBottom: 16 }}>
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSend}
              initialValues={{ fromAddress: selectedAddress }}
            >
              <Form.Item
                name="fromAddress"
                label="发送地址"
              >
                <Input disabled value={selectedAddress} />
              </Form.Item>
              
              <Form.Item
                name="toAddress"
                label="接收地址"
                rules={[{ required: true, message: '请输入接收地址' }]}
              >
                <Input placeholder="接收方钱包地址" />
              </Form.Item>
              
              <Form.Item
                name="amount"
                label="金额"
                rules={[
                  { required: true, message: '请输入金额' },
                  { type: 'number', min: 1, transform: val => Number(val), message: '金额必须大于0' }
                ]}
              >
                <Input placeholder="转账金额" type="number" />
              </Form.Item>
              
              <Form.Item>
                <Button 
                  type="primary" 
                  icon={<SendOutlined />}
                  htmlType="submit"
                  loading={sending}
                  disabled={!selectedAddress || balance <= 0}
                >
                  发送
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </TabPane>
        
        <TabPane 
          tab={
            <span>
              <WalletOutlined />
              交易记录
            </span>
          } 
          key="2"
        >
          <p>此功能尚未实现，将在后续版本中提供按地址查询交易记录的功能。</p>
        </TabPane>
      </Tabs>
    </div>
  );
};

export default Wallet; 