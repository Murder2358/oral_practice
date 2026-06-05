import { useEffect, useState } from 'react'
import { Card, Typography, Tag, Row, Col, message } from 'antd'
import {
  SolutionOutlined,
  CoffeeOutlined,
  TeamOutlined,
  SmileOutlined,
  GlobalOutlined,
  MessageOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { getScenes } from '../services/api'
import { useAppStore } from '../store'
import type { Scene } from '../types'

const { Title, Paragraph } = Typography

const iconMap: Record<string, React.ReactNode> = {
  briefcase: <SolutionOutlined style={{ fontSize: 32 }} />,
  coffee: <CoffeeOutlined style={{ fontSize: 32 }} />,
  team: <TeamOutlined style={{ fontSize: 32 }} />,
  smile: <SmileOutlined style={{ fontSize: 32 }} />,
  global: <GlobalOutlined style={{ fontSize: 32 }} />,
  message: <MessageOutlined style={{ fontSize: 32 }} />,
}

const difficultyColor: Record<string, string> = {
  beginner: 'green',
  intermediate: 'blue',
  advanced: 'red',
}

const difficultyLabel: Record<string, string> = {
  beginner: '初级',
  intermediate: '中级',
  advanced: '高级',
}

export default function Home() {
  const [scenes, setScenes] = useState<Scene[]>([])
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()
  const setSelectedScene = useAppStore((s) => s.setSelectedScene)

  useEffect(() => {
    getScenes()
      .then((res) => setScenes(res.data.scenes))
      .catch(() => message.error('Failed to load scenes'))
      .finally(() => setLoading(false))
  }, [])

  const handleSceneClick = (scene: Scene) => {
    setSelectedScene(scene)
    navigate(`/practice?scene=${scene.id}`)
  }

  return (
    <div style={{ padding: '40px 24px', maxWidth: 1200, margin: '0 auto' }}>
      <div style={{ textAlign: 'center', marginBottom: 48 }}>
        <Title level={1}>SpeakFlow</Title>
        <Paragraph style={{ fontSize: 18, color: '#666' }}>
          选择一个场景，开始你的英语口语练习
        </Paragraph>
      </div>

      <Row gutter={[24, 24]}>
        {scenes.map((scene) => (
          <Col xs={24} sm={12} lg={8} key={scene.id}>
            <Card
              hoverable
              loading={loading}
              onClick={() => handleSceneClick(scene)}
              style={{ height: '100%' }}
            >
              <div style={{ textAlign: 'center', marginBottom: 16 }}>
                {iconMap[scene.icon]}
              </div>
              <Title level={4} style={{ textAlign: 'center' }}>
                {scene.name}
              </Title>
              <Paragraph style={{ textAlign: 'center', color: '#666' }}>
                {scene.description}
              </Paragraph>
              <div style={{ textAlign: 'center' }}>
                {scene.difficulty.map((d) => (
                  <Tag key={d} color={difficultyColor[d]} style={{ margin: 2 }}>
                    {difficultyLabel[d]}
                  </Tag>
                ))}
              </div>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}
