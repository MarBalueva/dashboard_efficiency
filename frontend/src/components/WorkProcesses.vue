<template>
  <div class="page">
    <Sidebar />

    <main class="content">
      <header class="page-header">
        <h1>Рабочие процессы</h1>
        <p class="subtitle">Просмотр рабочих дней и метрик сотрудников</p>
      </header>

      <section class="table-section">
        <div class="table-header">
          <input
            v-model="search"
            placeholder="Поиск по сотруднику (ФИО)"
            class="search-input"
          />
        </div>

        <table class="employee-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>ФИО</th>
              <th>Начало дня</th>
              <th>Конец дня</th>
              <th>Звонки</th>
              <th>Выполненные задачи</th>
              <th>Work/Life Balance</th>
              <th>Удовлетворённость</th>
              <th>Продуктивность</th>
              <th>Действия</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="row in filteredData" :key="row.employee_id + row.start_work_day">
              <td>{{ row.employee_id }}</td>
              <td>{{ row.full_name }}</td>
              <td>{{ formatDate(row.start_work_day) }}</td>
              <td>{{ formatDate(row.end_work_day) }}</td>
              <td>{{ row.calls_count }}</td>
              <td>{{ row.completed_tasks }}</td>
              <td>{{ row.work_life_balance }}</td>
              <td>{{ row.satisfaction }}</td>
              <td>{{ row.productivity }}</td>
              <td class="actions">
                <button @click="remove(row.employee_id)">🗑️</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="loading" class="loading">Загрузка…</div>
        <div v-if="error" class="error">{{ error }}</div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import api from '../axios'

const data = ref([])
const loading = ref(false)
const error = ref('')
const search = ref('')

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get('/api/employees/work')
    data.value = res.data.map(r => ({
      ...r,
      start_work_day: r.start_work_day,
      end_work_day: r.end_work_day
    }))
  } catch {
    error.value = 'Ошибка загрузки данных'
  } finally {
    loading.value = false
  }
}

function remove(employee_id) {
  if (!confirm('Удалить записи сотрудника?')) return
  api.delete(`/api/employees/work/${employee_id}`).then(() => {
    data.value = data.value.filter(d => d.employee_id !== employee_id)
  }).catch(() => alert('Ошибка удаления'))
}

function formatDate(date) {
  return date ? new Date(date).toLocaleString() : '—'
}

const filteredData = computed(() => {
  const term = search.value.toLowerCase()
  return data.value.filter(row =>
    row.full_name.toLowerCase().includes(term)
  )
})

onMounted(fetchData)
</script>

<style scoped>
.page { display:flex; min-height:100vh; font-family: 'Nunito', sans-serif; }
.content { flex:1; padding: 36px; background: #f9fafb; box-sizing:border-box; }
.page-header h1 { margin:0 0 8px; font-size:24px; }
.subtitle { margin:0 0 18px; color:#64748B; }

.table-header { display:flex; justify-content:flex-start; align-items:center; margin-bottom:12px; }
.search-input {
  width:300px; padding:10px 14px; border-radius:8px;
  border:1px solid #cbd5e1; background-color:#ffffff; color:#1f2937;
  font-size:14px; box-shadow:0 1px 3px rgba(0,0,0,0.05);
}
.search-input::placeholder { color: #9ca3af; }

.employee-table { width:100%; border-collapse:collapse; text-align:left; background:#fff; }
.employee-table th, .employee-table td { padding:8px 12px; border-bottom:1px solid #e5e7eb; font-size:14px; }
.employee-table th { background:#f3f4f6; }
.employee-table th:hover { background:#e5e7eb; }

.actions button {
  background:#ffffff;
  border-radius:6px;
  cursor:pointer;
  font-size:16px;
  padding:4px 8px;
  transition:0.2s;
}
.actions button:hover { color:white; }

.loading { margin-top:12px; color:#64748B; }
.error { margin-top:12px; color:#EF4444; }
</style>
