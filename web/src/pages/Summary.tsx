import { Button, Typography, Card, Result } from 'antd'
import { useNavigate } from 'react-router-dom'

const { Title, Paragraph } = Typography

export default function Summary() {
  const navigate = useNavigate()

  return (
    <div style={{ padding: '40px 24px', maxWidth: 800, margin: '0 auto' }}>
      <Result
        status="success"
        title="Practice Complete!"
        subTitle="Your session summary will appear here after Day 3 implementation."
      />
      <Card style={{ marginTop: 24 }}>
        <Title level={4}>Coming Soon</Title>
        <Paragraph>课后总结功能将在 Day 3 开发，包括：</Paragraph>
        <ul>
          <li>对话记录回顾</li>
          <li>语法错误汇总与纠正</li>
          <li>发音评分详情</li>
          <li>能力雷达图</li>
          <li>练习进步趋势</li>
        </ul>
      </Card>
      <div style={{ textAlign: 'center', marginTop: 24 }}>
        <Button type="primary" size="large" onClick={() => navigate('/')}>
          Back to Home
        </Button>
      </div>
    </div>
  )
}
