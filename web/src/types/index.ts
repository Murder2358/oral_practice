export interface Scene {
  id: string
  name: string
  description: string
  icon: string
  difficulty: string[]
}

export interface Session {
  id: number
  scene_id: string
  difficulty: string
  score: number
  summary: string
  status: string
  created_at: string
  ended_at: string | null
}

export interface Message {
  id: number
  session_id: number
  role: 'user' | 'assistant' | 'system'
  content: string
  audio_url: string
  score: number
  created_at: string
}

export interface WSMessage {
  type: 'user' | 'assistant' | 'error' | 'score'
  content: string
  score?: number
}
