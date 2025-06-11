import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, Spin, message, Typography, Descriptions } from 'antd';
import { SaveOutlined, ReloadOutlined, SettingOutlined } from '@ant-design/icons';
import axios from 'axios';

const { Title } = Typography;

const NetworkConfig = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [config, setConfig] = useState({});
  const [isEditing, setIsEditing] = useState(false);

  const API_URL = 'http://localhost:8080/api/v1';

  // 获取当前网络配置
  const fetchConfig = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get(`${API_URL}/network/config`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      setConfig(response.data);
      form.setFieldsValue(response.data);
      setLoading(false);
    } catch (error) {
      console.error('获取网络配置失败:', error);
      message.error('获取网络配置失败: ' + (error.response?.data?.error || '请检查您的网络连接'));
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfig();
  }, []);

  // 更新网络配置
  const handleUpdate = async (values) => {
    setSaving(true);
    try {
      const token = localStorage.getItem('token');
      const response = await axios.put(`${API_URL}/network/config`, values, {
        headers: { Authorization: `Bearer ${token}` }
      });
      
      setConfig(response.data);
      message.success('网络配置已更新');
      setIsEditing(false);
      setSaving(false);
    } catch (error) {
      console.error('更新网络配置失败:', error);
      message.error('更新网络配置失败: ' + (error.response?.data?.error || '请检查参数是否有效'));
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
        <p>加载网络配置中...</p>
      </div>
    );
  }

  return (
    <div>
      <Title level={2}>区块链网络配置</Title>
      
      <Card 
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>网络参数</span>
            <div>
              <Button 
                type="default"
                icon={<ReloadOutlined />}
                onClick={fetchConfig}
                style={{ marginRight: 8 }}
              >
                刷新
              </Button>
              {!isEditing ? (
                <Button
                  type="primary"
                  icon={<SettingOutlined />}
                  onClick={() => setIsEditing(true)}
                >
                  编辑配置
                </Button>
              ) : (
                <Button
                  onClick={() => {
                    setIsEditing(false);
                    form.setFieldsValue(config);
                  }}
                >
                  取消
                </Button>
              )}
            </div>
          </div>
        }
      >
        {!isEditing ? (
          <Descriptions bordered column={1}>
            <Descriptions.Item label="区块生成时间间隔(秒)">
              {config.blockTime}
            </Descriptions.Item>
            <Descriptions.Item label="挖矿难度">
              {config.difficulty}
            </Descriptions.Item>
            <Descriptions.Item label="挖矿奖励">
              {config.miningReward}
            </Descriptions.Item>
            <Descriptions.Item label="交易费用">
              {config.txFee}
            </Descriptions.Item>
          </Descriptions>
        ) : (
          <Form
            form={form}
            layout="vertical"
            onFinish={handleUpdate}
            initialValues={config}
          >
            <Form.Item
              name="blockTime"
              label="区块生成时间间隔(秒)"
              rules={[
                { required: true, message: '请输入区块生成时间' },
                { type: 'number', min: 1, transform: val => Number(val), message: '时间必须大于0' }
              ]}
            >
              <Input type="number" placeholder="区块生成时间间隔" />
            </Form.Item>
            
            <Form.Item
              name="difficulty"
              label="挖矿难度"
              rules={[
                { required: true, message: '请输入挖矿难度' },
                { type: 'number', min: 1, transform: val => Number(val), message: '难度必须大于0' }
              ]}
            >
              <Input type="number" placeholder="挖矿难度" />
            </Form.Item>
            
            <Form.Item
              name="miningReward"
              label="挖矿奖励"
              rules={[
                { required: true, message: '请输入挖矿奖励' },
                { type: 'number', min: 0, transform: val => Number(val), message: '奖励必须大于等于0' }
              ]}
            >
              <Input type="number" placeholder="挖矿奖励" />
            </Form.Item>
            
            <Form.Item
              name="txFee"
              label="交易费用"
              rules={[
                { required: true, message: '请输入交易费用' },
                { type: 'number', min: 0, transform: val => Number(val), message: '费用必须大于等于0' }
              ]}
            >
              <Input type="number" placeholder="交易费用" />
            </Form.Item>
            
            <Form.Item>
              <Button 
                type="primary" 
                htmlType="submit" 
                icon={<SaveOutlined />}
                loading={saving}
              >
                保存配置
              </Button>
            </Form.Item>
          </Form>
        )}
      </Card>
      
      <div style={{ marginTop: 16 }}>
        <Card title="配置说明">
          <p><strong>区块生成时间间隔：</strong> 系统生成新区块的时间间隔，单位为秒。较短的时间会增加交易处理速度，但可能导致更多的分叉。</p>
          <p><strong>挖矿难度：</strong> 影响工作量证明算法的难度。数值越高，挖矿需要的计算资源越多。</p>
          <p><strong>挖矿奖励：</strong> 成功挖出新区块时矿工获得的奖励。</p>
          <p><strong>交易费用：</strong> 每笔交易需要支付的基础费用。</p>
        </Card>
      </div>
    </div>
  );
};

export default NetworkConfig; 