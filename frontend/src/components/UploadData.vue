<template>
  <div class="page">
    <Sidebar />

    <main class="content">
      <header class="page-header">
        <h1>Загрузить данные</h1>
        <p class="subtitle">Загрузите CSV/Excel с данными сотрудников для анализа</p>
      </header>

      <section class="upload-area" @drop.prevent="onDrop" @dragover.prevent>
  <div class="upload-inner">
    <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
      <path d="M11 16h2v-6h3l-4-5-4 5h3v6zM5 20h14v-2H5v2z" fill="#4F46E5"/>
    </svg>
    <p class="large">Перетащите файл сюда или</p>

    <!-- Контейнер для двух кнопок -->
    <div class="upload-buttons">
      <label class="upload-btn">
        <input type="file" ref="fileInput" @change="onFile" accept=".csv, .xlsx, .xls" />
        Выбрать файл
      </label>

      <a href="/template.xlsx" download class="btn-template">Скачать шаблон</a>
    </div>

    <p v-if="fileName" class="file-name">Выбран: {{ fileName }}</p>

    <button class="btn-indigo" :disabled="!file" @click="upload">Загрузить</button>

    <div v-if="message" :class="['message', messageType === 'error' ? 'error' : 'success']">{{ message }}</div>
  </div>
</section>

    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import api from '../axios' // axios instance

const file = ref(null)
const fileName = ref('')
const message = ref('')
const messageType = ref('')

function onFile(e) {
  const f = e.target.files[0]
  handleFile(f)
}

function onDrop(e) {
  const f = e.dataTransfer.files[0]
  handleFile(f)
}

function handleFile(f) {
  message.value = ''
  if (!f) return
  file.value = f
  fileName.value = f.name
}

async function upload() {
  if (!file.value) {
    messageType.value = 'error'
    message.value = 'Выберите файл'
    return
  }
  const form = new FormData()
  form.append('file', file.value)

  try {
    const res = await api.post('/api/upload', form, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    messageType.value = 'success'
    message.value = res.data?.message || 'Файл успешно загружен'
    file.value = null
    fileName.value = ''
  } catch (err) {
    messageType.value = 'error'
    message.value = err.response?.data?.message || 'Ошибка при загрузке'
  }
}
</script>

<style scoped>
.page { display:flex; min-height:100vh; font-family: 'Nunito', sans-serif; }
.content { flex:1; padding: 36px; background: var(--bg); box-sizing:border-box; }

/* header */
.page-header h1 { margin: 0 0 8px; font-size: 24px; }
.subtitle { margin:0 0 18px; color:#64748B; }

/* upload area */
.upload-area {
  display:flex; align-items:center; justify-content:center;
  width:100%; min-height:60vh;
}
.upload-inner {
  width:640px; border-radius:12px; background:#fff; padding:36px;
  display:flex; flex-direction:column; align-items:center; gap:14px;
  box-shadow: 0 8px 24px rgba(15,23,42,0.06);
  text-align:center;
}

.upload-buttons {
  display: flex;
  gap: 12px; /* расстояние между кнопками */
  margin-bottom: 12px; /* небольшой отступ снизу */
}
.upload-inner .large { font-size:18px; color:#475569; margin: 6px 0; }

.upload-btn {
  display: inline-block;
  padding: 10px 16px;
  border-radius: 8px;
  background: #EEF2FF;
  color: #4F46E5;
  font-weight: 600;
  cursor: pointer;
}

.upload-btn input { display:none; }

.btn-template {
  display: inline-block;
  padding: 10px 16px;
  border-radius: 8px;
  background: #F3F4F6;
  color: #1F2937;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-template:hover {
  background: #E5E7EB;
}

.file-name { color:#475569; font-size:14px; margin-top:6px; }
.message { margin-top:12px; font-size:14px; }
.message.error { color:#EF4444; }
.message.success { color:#10B981; }

/* reuse button style */
.btn-indigo { margin-top:12px; padding:10px 14px; background:#4F46E5; color:#fff; border:none; border-radius:8px; cursor:pointer; }
.btn-indigo:disabled { opacity:0.6; cursor:not-allowed; }
</style>
