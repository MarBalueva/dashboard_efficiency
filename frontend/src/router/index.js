import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import Register from '../components/Register.vue'
import Home from '../components/Home.vue'
import UploadData from '../components/UploadData.vue'


const routes = [
  { path: '/', name: 'home', component: Home },
  { path: '/login', name: 'login', component: Login },
  { path: '/register', name: 'register', component: Register },
  { path: '/upload-data', name: 'UploadData', component: UploadData },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
