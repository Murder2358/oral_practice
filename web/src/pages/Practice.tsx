import { useEffect, useRef, useState } from 'react'
import { Input, Button, Typography, Space, Select, message, Avatar, Spin, Tag } from 'antd'
import { SendOutlined, SoundOutlined, UserOutlined, RobotOutlined } from '@ant-design/icons'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { createSession, endSession } from '../services/api'
import { useAppStore } from '../store'
import { useWebSocket } from '../hooks/useWebSocket'
import type { WSMessage } from '../types'

const { Text } = Typography

interface ChatBubble {
  role: 'user' | 'assistant' | 'error'
  content: string
}

export default function Practice() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const sceneId = searchParams.get('scene') || 'free'
  const [difficulty, setDifficulty] = useState('intermediate')
  const [inputText, setInputText] = useState('')
  const [chatBubbles, setChatBubbles] = useState<ChatBubble[]>([])
  const [sessionId, setSessionId] = useState<number | null>(null)
  const [loading, setLoading] = useState(false)
  const chatEndRef = useRef<HTMLDivElement>(null)
  const selectedScene = useAppStore((s) => s.selectedScene)

  const handleWSMessage = (msg: WSMessage) => {
    setLoading(false)
    setChatBubbles((prev) => [...prev, { role: msg.type as 'user' | 'assistant' | 'error', content: msg.content }])
  }

  const { connected, connect, send, disconnect } = useWebSocket({
    onMessage: handleWSMessage,
    onOpen: () => message.success('AI partner connected'),
    onClose: () => {},
  })

  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [chatBubbles])

  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [])

  const handleStart = async () => {
    try {
      const res = await createSession({ scene_id: sceneId, difficulty })
      const sid = res.data.session.id
      setSessionId(sid)
      setChatBubbles([])
      connect(`/ws?session_id=${sid}&scene_id=${sceneId}&difficulty=${difficulty}`)
    } catch {
      message.error('Failed to create session')
    }
  }

  const handleEnd = async () => {
    disconnect()
    if (sessionId) {
      try {
        await endSession(sessionId)
      } catch {
        // ignore
      }
    }
    navigate(`/summary?session=${sessionId}`)
  }

  const handleSend = () => {
    const text = inputText.trim()
    if (!text) return

    setChatBubbles((prev) => [...prev, { role: 'user', content: text }])
    send({ type: 'user', content: text })
    setInputText('')
    setLoading(true)
  }

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column', background: '#f0f2f5' }}>
      {/* 顶部栏 */}
      <div
        style={{
          padding: '12px 24px',
          borderBottom: '1px solid #e8e8e8',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          background: '#fff',
          boxShadow: '0 1px 4px rgba(0,0,0,0.05)',
        }}
      >
        <Space>
          <Button onClick={() => { disconnect(); navigate('/') }}>Back</Button>
          <Text strong style={{ fontSize: 16 }}>{selectedScene?.name || 'Free Talk'}</Text>
          <Tag color={connected ? 'green' : 'default'}>{connected ? 'Connected' : 'Offline'}</Tag>
        </Space>
        <Space>
          <Select
            value={difficulty}
            onChange={setDifficulty}
            disabled={connected}
            style={{ width: 140 }}
            options={[
              { value: 'beginner', label: 'Beginner' },
              { value: 'intermediate', label: 'Intermediate' },
              { value: 'advanced', label: 'Advanced' },
            ]}
          />
          {!connected ? (
            <Button type="primary" onClick={handleStart}>
              Start Practice
            </Button>
          ) : (
            <Button danger onClick={handleEnd}>
              End Session
            </Button>
          )}
        </Space>
      </div>

      {/* 对话区 */}
      <div
        style={{
          flex: 1,
          overflow: 'auto',
          padding: '24px',
          display: 'flex',
          flexDirection: 'column',
          gap: 16,
        }}
      >
        {!connected && chatBubbles.length === 0 && (
          <div style={{ textAlign: 'center', marginTop: 80, color: '#999' }}>
            <RobotOutlined style={{ fontSize: 64, marginBottom: 16 }} />
            <div style={{ fontSize: 18 }}>Choose difficulty and click "Start Practice"</div>
          </div>
        )}

        {chatBubbles.map((bubble, idx) => (
          <div
            key={idx}
            style={{
              display: 'flex',
              justifyContent: bubble.role === 'user' ? 'flex-end' : 'flex-start',
            }}
          >
            <div
              style={{
                display: 'flex',
                alignItems: 'flex-start',
                gap: 8,
                maxWidth: '70%',
                flexDirection: bubble.role === 'user' ? 'row-reverse' : 'row',
              }}
            >
              <Avatar
                size={36}
                style={{
                  background: bubble.role === 'user' ? '#1677ff' : '#52c41a',
                  flexShrink: 0,
                }}
                icon={bubble.role === 'user' ? <UserOutlined /> : <SoundOutlined />}
              />
              <div
                style={{
                  padding: '10px 16px',
                  borderRadius: 12,
                  background: bubble.role === 'user' ? '#1677ff' : '#fff',
                  color: bubble.role === 'user' ? '#fff' : '#333',
                  boxShadow: '0 1px 4px rgba(0,0,0,0.08)',
                  wordBreak: 'break-word',
                  lineHeight: 1.6,
                }}
              >
                {bubble.content}
              </div>
            </div>
          </div>
        ))}

        {loading && (
          <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
            <Avatar size={36} style={{ background: '#52c41a' }} icon={<SoundOutlined />} />
            <div style={{ padding: '10px 16px', borderRadius: 12, background: '#fff', boxShadow: '0 1px 4px rgba(0,0,0,0.08)' }}>
              <Spin size="small" /> <Text type="secondary" style={{ marginLeft: 8 }}>AI is thinking...</Text>
            </div>
          </div>
        )}

        <div ref={chatEndRef} />
      </div>

      {/* 输入区 */}
      {connected && (
        <div
          style={{
            padding: '16px 24px',
            borderTop: '1px solid #e8e8e8',
            background: '#fff',
          }}
        >
          <Space.Compact style={{ width: '100%' }}>
            <Input
              size="large"
              placeholder="Type your message in English..."
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              onPressEnter={handleSend}
              disabled={loading}
              style={{ borderRadius: '8px 0 0 8px' }}
            />
            <Button
              type="primary"
              size="large"
              icon={<SendOutlined />}
              onClick={handleSend}
              disabled={loading || !inputText.trim()}
              style={{ borderRadius: '0 8px 8px 0', height: 40, width: 60 }}
            />
          </Space.Compact>
        </div>
      )}
    </div>
  )
}
