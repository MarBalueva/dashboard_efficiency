<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h2 style="margin-bottom: 8px;">Вход</h2>
      <p style="margin-bottom: 18px; color:#64748B;">
        Введите логин и пароль
      </p>

      <input
        v-model="login"
        :class="['input', errorField === 'login' || errorField === 'both' ? 'input--error' : '']"
        type="text"
        placeholder="Логин"
        @input="clearFieldError"
      />

      <input
        v-model="password"
        :class="['input', errorField === 'password' || errorField === 'both' ? 'input--error' : '']"
        type="password"
        placeholder="Пароль"
        @input="clearFieldError"
      />

      <button class="btn-indigo" @click="loginUser" :disabled="loading">
        {{ loading ? 'Вход...' : 'Войти' }}
      </button>

      <div
        v-if="serverMessage"
        class="error-text"
        :style="{ color: serverIsError ? '#EF4444' : '#10B981' }"
      >
        {{ serverMessage }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../axios'

const router = useRouter()

const login = ref('')
const password = ref('')
const loading = ref(false)

const errorField = ref('')
const serverMessage = ref('')
const serverIsError = ref(true)

function clearFieldError() {
  errorField.value = ''
  serverMessage.value = ''
}

function mapServerErrorToUI(code, message) {
  switch (code) {
    case 'user_not_found':
      errorField.value = 'login'
      serverMessage.value = message || 'Пользователь не найден'
      break
    case 'wrong_password':
      errorField.value = 'password'
      serverMessage.value = message || 'Неверный пароль'
      break
    case 'invalid_input':
      errorField.value = 'both'
      serverMessage.value = message || 'Проверьте введённые данные'
      break
    default:
      errorField.value = 'both'
      serverMessage.value = message || 'Ошибка авторизации'
  }
  serverIsError.value = true
}

async function loginUser() {
  errorField.value = ''
  serverMessage.value = ''

  if (!login.value || !password.value) {
    errorField.value = !login.value && !password.value ? 'both' : (!login.value ? 'login' : 'password')
    serverMessage.value = 'Заполните все поля'
    serverIsError.value = true
    return
  }

  loading.value = true
  try {
    const res = await api.post('/auth/login', {
      login: login.value,
      password: password.value
    })

    const token = res.data.token
    if (!token) {
      serverMessage.value = 'Сервер не вернул токен'
      serverIsError.value = true
      return
    }

    localStorage.setItem('token', token)
    api.defaults.headers.common.Authorization = `Bearer ${token}`

    serverMessage.value = 'Успешный вход'
    serverIsError.value = false

    setTimeout(() => {
      router.push('/dashboard')
    }, 300)
  } catch (err) {
    const resp = err.response
    if (resp?.data) {
      mapServerErrorToUI(resp.data.error, resp.data.message)
    } else {
      errorField.value = 'both'
      serverMessage.value = 'Ошибка сети или сервера'
      serverIsError.value = true
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.input--error {
  border-color: #EF4444 !important;
  background: #FFF1F2;
}

.error-text {
  margin-top: 12px;
  font-size: 14px;
}
</style>
