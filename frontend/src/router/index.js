import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import Register from '../components/Register.vue'
import Home from '../components/Home.vue'
import UploadData from '../components/UploadData.vue'
import EmployeeList from '../components/EmployeeList.vue'
import Profile from '../components/Profile.vue'
import Dashboard from '../components/Dashboard.vue'
import DictPage from '../components/DictPage.vue'
import WorkProcesses from '../components/WorkProcesses.vue'


const routes = [
  { path: '/', name: 'home', component: Home },
  { path: '/login', name: 'login', component: Login },
  { path: '/register', name: 'register', component: Register },
  { path: '/upload-data', name: 'UploadData', component: UploadData },
  { path: '/employees', name: 'EmployeeList', component: EmployeeList},
  { path: '/profile', name: 'Profile', component: Profile },
  { path: '/dashboard', name: 'Dashboard', component: Dashboard },
  { path: '/dict', name: 'DictPage', component: DictPage },
  { path: '/employees/work', name: 'WorkProcesses', component: WorkProcesses },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
