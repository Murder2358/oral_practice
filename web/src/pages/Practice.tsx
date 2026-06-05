import { useEffect, useRef, useState } from 'react'
import { Input, Button, Typography, Card, Space, Select, message, Avatar } from 'antd'
import { SendOutlined, SoundOutlined, UserOutlined } from '@ant-design/icons'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { useAppStore } from '../store'
import type { WSMessage } from '../types'

const { Text } = Typography

export default function Practice() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const sceneId = searchParams.get('scene') || 'free'
  const [difficulty, setDifficulty] = useState('intermediate')
  const [inputText, setInputText] = useState('')
  const [chatMessages, setChatMessages] = useState<WSMessage[]>([])
  const [wsConnected, setWsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)
  const chatEndRef = useRef<HTMLDivElement>(null)
  const selectedScene = useAppStore((s) => s.selectedScene)

  const connectWS = (diff: string) => {
    const ws = new WebSocket(
      `ws://localhost:8080/ws?scene_id=${sceneId}&difficulty=${diff}`
    )
    wsRef.current = ws

    ws.onopen = () => {
      setWsConnected(true)
      message.success('Connected to AI partner')
    }

    ws.onmessage = (event) => {
      const msg: WSMessage = JSON.parse(event.data)
      setChatMessages((prev) => [...prev, msg])
    }

    ws.onclose = () => {
      setWsConnected(false)
    }

    ws.onerror = () => {
      message.error('Connection failed')
    }
  }

  useEffect(() => {
    return () => {
      wsRef.current?.close()
    }
  }, [])

  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [chatMessages])

  const handleStart = () => {
    setChatMessages([])
    connectWS(difficulty)
  }

  const handleSend = () => {
    if (!inputText.trim() || !wsRef.current) return
    const msg = { type: 'user' as const, content: inputText.trim() }
    setChatMessages((prev) => [...prev, msg])
    wsRef.current.send(JSON.stringify(msg))
    setInputText('')
  }

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      {/* 顶部栏 */}
      <div
        style={{
          padding: '12px 24px',
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          background: '#fff',
        }}
      >
        <Space>
          <Button onClick={() => navigate('/')}>Back</Button>
          <Text strong>{selectedScene?.name || 'Free Talk'}</Text>
          <Select
            value={difficulty}
            onChange={setDifficulty}
            style={{ width: 120 }}
            options={[
              { value: 'beginner', label: 'Beginner' },
              { value: 'intermediate', label: 'Intermediate' },
              { value: 'advanced', label: 'Advanced' },
            ]}
          />
        </Space>
        <Space>
          <Text type={wsConnected ? 'success' : 'secondary'}>
            {wsConnected ? 'Connected' : 'Disconnected'}
          </Text>
          {!wsConnected ? (
            <Button type="primary" onClick={handleStart}>
              Start Practice
            </Button>
          ) : (
            <Button
              danger
              onClick={() => {
                wsRef.current?.close()
                navigate('/summary')
              }}
            >
              End Session
            </Button>
          )}
        </Space>
      </div>

      {/* 对话区 */}
      <div style={{ flex: 1, overflow: 'auto', padding: '24px', background: '#f5f5f5' }}>
        {chatMessages.map((msg, idx) => (
          <div
            key={idx}
            style={{
              display: 'flex',
              justifyContent: msg.type === 'user' ? 'flex-end' : 'flex-start',
              marginBottom: 16,
            }}
          >
            <Card
              size="small"
              style={{
                maxWidth: '70%',
                background: msg.type === 'user' ? '#1677ff' : '#fff',
                color: msg.type === 'user' ? '#fff' : '#000',
              }}
            >
              <Space>
                {msg.type !== 'user' && (
                  <Avatar size="small" icon={<SoundOutlined />} />
                )}
                <Text
                  style={{
                    color: msg.type === 'user' ? '#fff' : '#000',
                  }}
                >
                  {msg.content}
                </Text>
                {msg.type === 'user' && (
                  <Avatar size="small" icon={<UserOutlined />} />
                )}
              </Space>
            </Card>
          </div>
        ))}
        <div ref={chatEndRef} />
      </div>

      {/* 输入区 */}
      {wsConnected && (
        <div
          style={{
            padding: '12px 24px',
            borderTop: '1px solid #f0f0f0',
            background: '#fff',
            display: 'flex',
            gap: 8,
          }}
        >
          <Input
            size="large"
            placeholder="Type your message..."
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            onPressEnter={handleSend}
          />
          <Button
            type="primary"
            size="large"
            icon={<SendOutlined />}
            onClick={handleSend}
          />
        </div>
      )}
    </div>
  )
}
