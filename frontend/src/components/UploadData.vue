<template>
  <div class="page">
    <Sidebar />

    <main class="content">
      <header class="page-header">
        <h1>Загрузить данные</h1>
        <p class="subtitle">Загрузите CSV/Excel с данными сотрудников для анализа</p>
      </header>

      <section class="upload-area">
        <!-- Загрузка файла -->
        <div class="upload-inner" v-if="!previewData.length">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
            <path d="M11 16h2v-6h3l-4-5-4 5h3v6zM5 20h14v-2H5v2z" fill="#4F46E5"/>
          </svg>
          <p class="large">Перетащите файл сюда или</p>

          <div class="upload-buttons">
            <label class="upload-btn">
              <input type="file" ref="fileInput" @change="onFile" accept=".csv, .xlsx, .xls" />
              Выбрать файл
            </label>
            <a href="/template.xlsx" download class="btn-template">Скачать шаблон</a>
          </div>

          <p v-if="fileName" class="file-name">Выбран: {{ fileName }}</p>

          <button class="btn-indigo" :disabled="!file" @click="upload">Загрузить</button>

          <div v-if="message" :class="['message', messageType === 'error' ? 'error' : 'success']">
            {{ message }}
          </div>
        </div>

        <!-- Предпросмотр данных -->
        <div v-else class="upload-inner preview-inner">
          <h2>Предпросмотр загруженных данных</h2>
          <div class="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>EmployeeID</th>
                  <th>StartWorkDay</th>
                  <th>EndWorkDay</th>
                  <th>CallsCount</th>
                  <th>CompletedTasks</th>
                  <th>WorkLifeBalance</th>
                  <th>Satisfaction</th>
                  <th>Productivity</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, idx) in previewData" :key="idx">
                  <td>{{ row.employee_id }}</td>
                  <td>{{ formatDateTime(row.start_work_day) }}</td>
                  <td>{{ formatDateTime(row.end_work_day) }}</td>
                  <td>{{ row.calls_count }}</td>
                  <td>{{ row.completed_tasks }}</td>
                  <td>{{ row.work_life_balance }}</td>
                  <td>{{ row.satisfaction }}</td>
                  <td>{{ row.productivity }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="preview-buttons">
            <button class="btn-gray" @click="cancelPreview">Отмена</button>
            <button class="btn-indigo" @click="confirmUpload">Подтвердить</button>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import api from '../axios'

const router = useRouter()

const file = ref(null)
const fileName = ref('')
const message = ref('')
const messageType = ref('')
const previewData = ref([])

function handleFile(f) {
  message.value = ''
  if (!f) return
  file.value = f
  fileName.value = f.name
}

function onFile(e) {
  handleFile(e.target.files[0])
}

async function upload() {
  if (!file.value) return

  const form = new FormData()
  form.append('file', file.value)

  try {
    const res = await api.post('/api/upload', form, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    if (!res.data.preview || res.data.preview.length === 0) {
      messageType.value = 'error'
      message.value = 'Файл не содержит данных'
      return
    }

    previewData.value = res.data.preview
  } catch (err) {
    messageType.value = 'error'
    message.value = err.response?.data?.message || 'Ошибка при загрузке'
  }
}

function cancelPreview() {
  previewData.value = []
  file.value = null
  fileName.value = ''
  message.value = ''
  messageType.value = ''
}

async function confirmUpload() {
  try {
    await api.post('/api/upload/confirm', previewData.value)
    alert('Данные успешно сохранены')
    router.push('/employees/work')
  } catch (err) {
    alert(err.response?.data?.message || 'Ошибка при сохранении')
  }
}

function formatDateTime(date) {
  if (!date) return '—'
  const d = new Date(date)
  return d.toLocaleDateString() + ' ' + d.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
}
</script>

<style scoped>
.upload-inner { 
  width: 700px; 
  border-radius: 12px; 
  background: #fff; 
  padding: 36px; 
  display: flex; 
  flex-direction: column; 
  align-items: center; 
  gap: 14px; 
  box-shadow: 0 8px 24px rgba(15,23,42,0.06); 
  text-align: center; 
}

.upload-buttons {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.upload-btn {
  background: #EEF2FF;
  color: #4F46E5;
  position: relative;
  overflow: hidden;
  display: inline-block;
  padding: 10px 16px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

.upload-btn input {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.btn-template {
  display: inline-block;
  background: #F3F4F6;
  color: #1F2937;
  padding: 10px 16px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  cursor: pointer;
}

.btn-template:hover { background:#E5E7EB; }

.preview-inner {
  display: flex;
  flex-direction: column;
  align-items: stretch; 
  gap: 16px;
  text-align: center;
}

.preview-buttons {
  display: flex;
  justify-content: center; 
  gap: 16px;
}

.preview-buttons button {
  width: 150px;
  height: 42px;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  justify-content: center; 
  align-items: center;
  text-align: center;
}

.btn-indigo {
  background: #4F46E5;
  color: #fff;
  border: none;
}

.btn-gray {
  background: #F3F4F6;
  color: #1F2937;
  border: none;
}

.table-wrapper {
  width: 100%;
  overflow-x: auto;
  border: 1px solid #ddd;
  margin-bottom: 16px;
}

.table-wrapper table {
  border-collapse: collapse;
  min-width: 900px;
}

th, td {
  border: 1px solid #ddd;
  padding: 6px 8px;
  text-align: left;
  white-space: nowrap;
}
</style>
