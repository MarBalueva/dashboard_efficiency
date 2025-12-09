<template>
  <aside class="sidebar">
    <!-- Верхний логотип -->
    <div class="sidebar-top">
      <div class="company-icon">
        <svg width="40" height="40" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" fill="#4F46E5"/>
          <text x="50%" y="53%" text-anchor="middle" fill="#fff" font-size="10">C</text>
        </svg>
      </div>
    </div>

    <!-- Среднее меню -->
    <nav class="sidebar-middle">
      <router-link
        to="/dashboard"
        class="icon-btn"
        :class="{ active: isActive('/dashboard') }"
      >
        <svg width="25" height="25" viewBox="0 0 24 24" fill="none">
          <path d="M3 13h8V3H3v10zM13 21h8V11h-8v10zM13 3v6h8V3h-8zM3 21h8v-6H3v6z"
                fill="currentColor"/>
        </svg>
      </router-link>

      <router-link
        to="/upload-data"
        class="icon-btn"
        :class="{ active: isActive('/upload-data') }"
      >
        <svg width="25" height="25" viewBox="0 0 24 24" fill="none">
          <path d="M11 16h2V8h3l-4-5-4 5h3v8zM5 20h14v-2H5v2z"
                fill="currentColor"/>
        </svg>
      </router-link>

      <!-- Новый пункт: список сотрудников -->
      <router-link
        to="/employees"
        class="icon-btn"
        :class="{ active: isActive('/employees') }"
      >
        <svg width="25" height="25" viewBox="0 0 24 24" fill="none">
          <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 2.07 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z" fill="currentColor"/>
        </svg>
      </router-link>
    </nav>

    <!-- Нижняя часть -->
    <div class="sidebar-bottom">
      <router-link to="/profile" class="icon-btn" :class="{ active: isActive('/profile') }">
        <svg width="25" height="25" viewBox="0 0 24 24" fill="none">
          <path d="M12 12c2.7 0 5-2.3 5-5s-2.3-5-5-5-5 2.3-5 5 2.3 5 5 5zm0 2c-4 0-8 2-8 6v2h16v-2c0-4-4-6-8-6z" fill="currentColor"/>
        </svg>
      </router-link>

      <button class="icon-btn logout" @click="logout">
        <svg width="25" height="25" viewBox="0 0 24 24" fill="none">
          <path d="M16 13v-2H7V8l-5 4 5 4v-3zM20 3h-8v2h8v14h-8v2h8c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2z" fill="currentColor"/>
        </svg>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

function isActive(path) {
  return route.path === path
}

function logout() {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<style scoped>
/* Сайдбар */
.sidebar {
  width: 80px;
  height: 90vh;
  background: #ffffff;
  position: fixed;
  left: 25px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 20px 0;
  box-shadow: 2px 0 8px rgba(15,23,42,0.05);
  border-radius: 20px;
}

.sidebar-middle, .sidebar-bottom {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
}

.icon-btn {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: transparent;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: #475569;
}

.icon-btn.active {
  background: #EEF2FF;
  color: #4F46E5;
}

.logout {
  color: #EF4444;
}

.logout:hover {
  background: rgba(239,68,68,0.06);
}
</style>
