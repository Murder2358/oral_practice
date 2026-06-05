import { create } from 'zustand'
import type { Scene, Session, Message } from '../types'

interface AppState {
  // 场景
  scenes: Scene[]
  setScenes: (scenes: Scene[]) => void

  // 当前会话
  currentSession: Session | null
  setCurrentSession: (session: Session | null) => void

  // 对话消息
  messages: Message[]
  addMessage: (message: Message) => void
  clearMessages: () => void

  // WebSocket 连接状态
  wsConnected: boolean
  setWsConnected: (connected: boolean) => void

  // 当前场景
  selectedScene: Scene | null
  setSelectedScene: (scene: Scene | null) => void
}

export const useAppStore = create<AppState>((set) => ({
  scenes: [],
  setScenes: (scenes) => set({ scenes }),

  currentSession: null,
  setCurrentSession: (session) => set({ currentSession: session }),

  messages: [],
  addMessage: (message) => set((state) => ({ messages: [...state.messages, message] })),
  clearMessages: () => set({ messages: [] }),

  wsConnected: false,
  setWsConnected: (connected) => set({ wsConnected: connected }),

  selectedScene: null,
  setSelectedScene: (scene) => set({ selectedScene: scene }),
}))
