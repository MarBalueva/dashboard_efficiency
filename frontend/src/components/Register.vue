<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h2 style="margin-bottom: 10px;">Регистрация</h2>

      <input
        v-model="email"
        :class="['input', errorField === 'email' || errorField === 'both' ? 'input--error' : '']"
        type="email"
        placeholder="Email"
        @input="clearFieldError"
      />
      <input
        v-model="name"
        :class="['input', errorField === 'name' || errorField === 'both' ? 'input--error' : '']"
        type="text"
        placeholder="Имя"
        @input="clearFieldError"
      />
      <input
        v-model="password"
        :class="['input', errorField === 'password' || errorField === 'both' ? 'input--error' : '']"
        type="password"
        placeholder="Пароль (минимум 5 символов)"
        @input="clearFieldError"
      />

      <button class="btn-indigo" @click="register" :disabled="loading">
        {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
      </button>

      <div class="hint" style="margin-top:18px;">
        Уже есть аккаунт?
        <router-link to="/login">Войти</router-link>
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
const name = ref('')
const password = ref('')
const loading = ref(false)

const errorField = ref('')
const serverMessage = ref('')
const serverIsError = ref(true)

function clearFieldError() {
  errorField.value = ''
  serverMessage.value = ''
}

function mapRegisterError(code, msg) {
  switch (code) {
    case 'email_exists':
      errorField.value = 'email'
      serverMessage.value = msg || 'Пользователь с таким email уже существует'
      break
    case 'invalid_input':
      errorField.value = 'both'
      serverMessage.value = msg || 'Проверьте вводимые данные'
      break
    default:
      errorField.value = 'both'
      serverMessage.value = msg || 'Ошибка регистрации'
  }
  serverIsError.value = true
}

function validate() {
  if (!email.value || !name.value || !password.value) {
    errorField.value = !email.value && !name.value && !password.value ? 'both'
      : (!email.value ? 'email' : (!name.value ? 'name' : 'password'))
    serverMessage.value = 'Заполните все поля'
    serverIsError.value = true
    return false
  }
  if (password.value.length < 5) {
    errorField.value = 'password'
    serverMessage.value = 'Пароль должен быть не менее 5 символов'
    serverIsError.value = true
    return false
  }
  return true
}

async function register() {
  serverMessage.value = ''
  errorField.value = ''
  if (!validate()) return

  loading.value = true
  try {
    await api.post('/auth/register', {
      email: email.value,
      password: password.value,
      name: name.value
    })

    serverMessage.value = 'Пользователь создан. Перенаправление на вход...'
    serverIsError.value = false
    setTimeout(() => router.push('/login'), 800)
  } catch (err) {
    const resp = err.response
    if (resp && resp.data) {
      mapRegisterError(resp.data.error, resp.data.message || resp.data.error)
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
