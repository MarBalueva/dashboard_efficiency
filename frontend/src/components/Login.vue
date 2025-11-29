<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h2 style="margin-bottom: 8px;">Вход</h2>
      <p style="margin-bottom: 18px; color:#64748B;">Введите данные для входа</p>

      <input
        v-model="email"
        :class="['input', errorField === 'email' || errorField === 'both' ? 'input--error' : '']"
        type="email"
        placeholder="Email"
        @input="clearFieldError"
      />

      <input
        v-model="password"
        :class="['input', errorField === 'password' || errorField === 'both' ? 'input--error' : '']"
        type="password"
        placeholder="Пароль"
        @input="clearFieldError"
      />

      <button class="btn-indigo" @click="login" :disabled="loading">
        {{ loading ? 'Вход...' : 'Войти' }}
      </button>

      <div class="hint" style="margin-top:18px;">
        Нет аккаунта?
        <router-link to="/register">Зарегистрироваться</router-link>
      </div>

      <div v-if="serverMessage" class="error-text" :style="{ color: serverIsError ? '#EF4444' : '#10B981' }">
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
const email = ref('')
const password = ref('')
const loading = ref(false)

const errorField = ref('')        // '', 'email', 'password', 'both'
const serverMessage = ref('')     // текст для показа
const serverIsError = ref(true)   // true -> red, false -> success

function clearFieldError() {
  errorField.value = ''
  serverMessage.value = ''
}

function mapServerErrorToUI(errCode, msg) {
  // errCode ожидается как короткий код, например: "user_not_found", "wrong_password", "invalid_input", "email_exists"
  if (!errCode) {
    serverMessage.value = msg || 'Ошибка сервера, попробуйте позже'
    serverIsError.value = true
    return
  }

  switch (errCode) {
    case 'user_not_found':
      errorField.value = 'email'
      serverMessage.value = msg || 'Пользователь с таким email не найден'
      break
    case 'wrong_password':
      errorField.value = 'password'
      serverMessage.value = msg || 'Неверный пароль'
      break
    case 'invalid_input':
      // message может содержать подсказку
      errorField.value = 'both'
      serverMessage.value = msg || 'Проверьте введённые данные'
      break
    default:
      // на всякий случай: если сервис вернул общий код
      errorField.value = 'both'
      serverMessage.value = msg || 'Неверный email или пароль'
  }
  serverIsError.value = true
}

async function login() {
  serverMessage.value = ''
  errorField.value = ''

  if (!email.value || !password.value) {
    errorField.value = !email.value && !password.value ? 'both' : (!email.value ? 'email' : 'password')
    serverMessage.value = 'Заполните все поля'
    serverIsError.value = true
    return
  }

  loading.value = true
  try {
    const res = await api.post('/auth/login', {
      email: email.value,
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

    // краткая пауза, чтобы пользователь увидел сообщение
    setTimeout(() => router.push('/'), 400)
  } catch (err) {
    // ожидаем формат: { error: "code", message: "текст" }
    const resp = err.response
    if (resp && resp.data) {
      const code = resp.data.error
      const msg = resp.data.message || resp.data.error || resp.data
      mapServerErrorToUI(code, msg)
    } else {
      serverMessage.value = 'Ошибка сети или сервера'
      serverIsError.value = true
      errorField.value = 'both'
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
