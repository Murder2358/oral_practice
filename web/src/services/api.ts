import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

export interface CreateSessionReq {
  scene_id: string
  difficulty: string
}

export const getScenes = () => api.get('/scenes')

export const createSession = (data: CreateSessionReq) => api.post('/sessions', data)

export const getSession = (id: number) => api.get(`/sessions/${id}`)

export default api
