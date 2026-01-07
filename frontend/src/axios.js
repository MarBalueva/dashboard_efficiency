import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
  timeout: 10000,
})

const token = localStorage.getItem('token')
if (token) {
  api.defaults.headers.common.Authorization = `Bearer ${token}`
}

export default api
