<template>
  <div class="page">
    <Sidebar />

    <main class="content">
      <header class="page-header">
        <h1>Список сотрудников</h1>
        <p class="subtitle">Текущие загруженные данные сотрудников</p>
      </header>

      <section class="table-section">
        <input
          type="text"
          v-model="search"
          placeholder="Поиск по Employee ID или Department"
          class="search-input"
        />

        <div class="table-wrapper">
          <table class="employee-table">
            <thead>
              <tr>
                <th @click="sortBy('employee_id')">ID</th>
                <th @click="sortBy('age')">Возраст</th>
                <th @click="sortBy('department')">Отдел</th>
                <th @click="sortBy('job_level')">Уровень</th>
                <th @click="sortBy('years_at_company')">Лет в компании</th>
                <th @click="sortBy('monthly_hours_worked')">Часы/мес</th>
                <th>Удалённая работа</th>
                <th @click="sortBy('meetings_per_week')">Встреч/нед</th>
                <th @click="sortBy('tasks_completed_per_day')">Задач/день</th>
                <th @click="sortBy('overtime_hours_per_week')">Сверхчасы</th>
                <th>Баланс работы/жизни</th>
                <th @click="sortBy('job_satisfaction')">Удовлетворённость</th>
                <th @click="sortBy('productivity_score')">Продуктивность</th>
                <th @click="sortBy('annual_salary')">Зарплата</th>
                <th @click="sortBy('absences_per_year')">Прогулы</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="emp in filteredEmployees"
                :key="emp.id"
              >
                <td>{{ emp.employee_id }}</td>
                <td>{{ emp.age }}</td>
                <td>{{ emp.department }}</td>
                <td>{{ emp.job_level }}</td>
                <td>{{ emp.years_at_company }}</td>
                <td>{{ emp.monthly_hours_worked }}</td>
                <td>{{ emp.remote_work ? 'Да' : 'Нет' }}</td>
                <td>{{ emp.meetings_per_week }}</td>
                <td>{{ emp.tasks_completed_per_day }}</td>
                <td>{{ emp.overtime_hours_per_week }}</td>
                <td>{{ emp.work_life_balance }}</td>
                <td>{{ emp.job_satisfaction }}</td>
                <td>{{ emp.productivity_score }}</td>
                <td>{{ emp.annual_salary }}</td>
                <td>{{ emp.absences_per_year }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="loading" class="loading">Загрузка...</div>
        <div v-if="error" class="error">{{ error }}</div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import api from '../axios'

const employees = ref([])
const loading = ref(false)
const error = ref('')
const search = ref('')
const sortKey = ref('')
const sortAsc = ref(true)

async function fetchEmployees() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get('/api/employees')
    employees.value = res.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Ошибка при загрузке данных'
  } finally {
    loading.value = false
  }
}

function sortBy(key) {
  if (sortKey.value === key) {
    sortAsc.value = !sortAsc.value
  } else {
    sortKey.value = key
    sortAsc.value = true
  }
}

const filteredEmployees = computed(() => {
  let list = employees.value

  if (search.value) {
    const term = search.value.toLowerCase()
    list = list.filter(
      e =>
        String(e.employee_id).toLowerCase().includes(term) ||
        (e.department && e.department.toLowerCase().includes(term))
    )
  }

  if (sortKey.value) {
    list = [...list].sort((a, b) => {
      const valA = a[sortKey.value]
      const valB = b[sortKey.value]
      if (valA == null) return 1
      if (valB == null) return -1
      if (valA === valB) return 0
      return sortAsc.value
        ? valA > valB ? 1 : -1
        : valA < valB ? 1 : -1
    })
  }

  return list
})

onMounted(fetchEmployees)
</script>

<style scoped>
.page { display:flex; min-height:100vh; font-family: 'Nunito', sans-serif; }
.content { flex:1; padding: 36px; background: var(--bg); box-sizing:border-box; }

/* header */
.page-header h1 { margin: 0 0 8px; font-size: 24px; }
.subtitle { margin:0 0 18px; color:#64748B; }

/* таблица */
.table-section { margin-top: 20px; }
.table-wrapper {
  overflow-x: auto;
  max-height: 60vh;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
}
.employee-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}
.employee-table th, .employee-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #e5e7eb;
  font-size: 14px;
}
.employee-table th {
  background: #f3f4f6;
  cursor: pointer;
}
.employee-table th:hover {
  background: #e5e7eb;
}

.search-input {
  width: 100%;
  max-width: 300px;
  padding: 10px 14px;       /* увеличенный padding для высоты */
  margin-bottom: 12px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background-color: #ffffff; /* белый фон */
  color: #1f2937;            /* тёмный текст */
  font-size: 14px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1); /* лёгкая тень */
}

.search-input::placeholder {
  color: #94a3b8;            /* серый для плейсхолдера */
}

.search-input:focus {
  outline: none;
  border-color: #4F46E5; /* синий акцент при фокусе */
  box-shadow: 0 0 0 2px rgba(79,70,229,0.2);
}

.loading { margin-top: 12px; color: #64748B; }
.error { margin-top: 12px; color: #EF4444; }
</style>
