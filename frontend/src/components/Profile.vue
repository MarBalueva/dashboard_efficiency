<template>
  <div class="page">
    <Sidebar />

    <main class="content">
      <header class="page-header">
        <h1>Профиль пользователя</h1>
        <p class="subtitle">Вы можете изменить ФИО и пароль</p>
      </header>

      <section class="profile-section">
        <div class="profile-form">
          <!-- Логин (readonly) -->
          <div class="profile-field">
            <label>Логин</label>
            <input type="text" :value="profile.login" readonly />
          </div>

          <div class="profile-field">
            <label>Фамилия</label>
            <input
              type="text"
              v-model="form.last_name"
              placeholder="Введите фамилию"
            />
          </div>

          <div class="profile-field">
            <label>Имя</label>
            <input
              type="text"
              v-model="form.first_name"
              placeholder="Введите имя"
            />
          </div>

          <div class="profile-field">
            <label>Отчество</label>
            <input
              type="text"
              v-model="form.middle_name"
              placeholder="Введите отчество"
            />
          </div>

          <div class="profile-field">
            <label>Новый пароль</label>
            <input
              type="password"
              v-model="form.password"
              placeholder="Введите новый пароль"
            />
          </div>

          <button class="btn-indigo" @click="updateProfile" :disabled="loading">
            {{ loading ? 'Сохраняем...' : 'Сохранить изменения' }}
          </button>

          <div
            v-if="message"
            :class="['message', error ? 'error' : 'success']"
          >
            {{ message }}
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import api from '../axios'

const profile = ref({
  login: '',
  employee: {
    last_name: '',
    first_name: '',
    middle_name: ''
  }
})

const form = ref({
  last_name: '',
  first_name: '',
  middle_name: '',
  password: ''
})

const loading = ref(false)
const message = ref('')
const error = ref(false)

async function fetchProfile() {
  loading.value = true
  error.value = false
  try {
    const res = await api.get('/api/profile')
    profile.value = res.data

    form.value.last_name = res.data.employee.last_name
    form.value.first_name = res.data.employee.first_name
    form.value.middle_name = res.data.employee.middle_name
  } catch (err) {
    message.value = err.response?.data?.error || 'Ошибка загрузки профиля'
    error.value = true
  } finally {
    loading.value = false
  }
}

async function updateProfile() {
  loading.value = true
  message.value = ''
  error.value = false

  try {
    const payload = {
      last_name: form.value.last_name,
      first_name: form.value.first_name,
      middle_name: form.value.middle_name
    }

    if (form.value.password) {
      payload.password = form.value.password
    }

    const res = await api.put('/api/profile', payload)

    message.value = res.data.message || 'Профиль обновлён'

    profile.value.employee.last_name = form.value.last_name
    profile.value.employee.first_name = form.value.first_name
    profile.value.employee.middle_name = form.value.middle_name

    form.value.password = ''
  } catch (err) {
    message.value = err.response?.data?.error || 'Ошибка при обновлении профиля'
    error.value = true
  } finally {
    loading.value = false
  }
}

onMounted(fetchProfile)
</script>

<style scoped>
.content { padding: 36px; flex:1; }
.page-header h1 { margin-bottom: 8px; }
.subtitle { margin-bottom: 18px; color:#64748B; }

.profile-section {
  display: flex;
  justify-content: center; 
  padding: 20px;
}

.profile-form {
  width: 100%;
  max-width: 600px;
  background:#fff;
  padding: 36px;
  margin: 0 auto;
  border-radius:12px;
  box-shadow:0 8px 24px rgba(15,23,42,0.06);
  display:flex;
  flex-direction:column;
  gap:16px;
}

.profile-field label {
  font-weight:600;
  margin-bottom:4px;
  display:block;
}

.profile-field input {
  width: 100%;
  padding: 10px 14px; 
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background: #fff;
  color: #000;
  font-size: 14px;
  box-sizing: border-box;
}

.profile-field input::placeholder {
  color: #94a3b8;      
  opacity: 1;           
}

.profile-field input:focus {
  outline: none;
  border-color: #4F46E5;
  box-shadow: 0 0 0 2px rgba(79,70,229,0.2);
}


.profile-field input[readonly] {
  background: #f3f4f6; 
  cursor: not-allowed;
}

.btn-indigo {
  padding:12px 16px;
  background:#4F46E5;
  color:#fff;
  border:none;
  border-radius:8px;
  font-weight:600;
  cursor:pointer;
}

.message.success { color: #16a34a; }
.message.error { color: #ef4444; }
</style>
