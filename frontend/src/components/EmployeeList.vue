<template>
  <div class="page">
    <Sidebar />

    <main class="content">
      <header class="page-header">
        <h1>Сотрудники</h1>
        <p class="subtitle">Кадровые данные</p>
      </header>

      <section class="table-section">
        <!-- Поиск по ФИО -->
        <div class="table-header">
          <input
            v-model="search"
            placeholder="Поиск по ФИО"
            class="search-input"
          />
          <button class="add-btn" @click="openAddModal">Добавить сотрудника</button>
        </div>

        <!-- Таблица -->
        <table class="employee-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>
                Фамилия
                <select v-model="filters.last_name" class="column-filter">
                  <option value="">Все</option>
                  <option v-for="name in uniqueValues('last_name')" :key="name">{{ name }}</option>
                </select>
              </th>
              <th>
                Имя
                <select v-model="filters.first_name" class="column-filter">
                  <option value="">Все</option>
                  <option v-for="name in uniqueValues('first_name')" :key="name">{{ name }}</option>
                </select>
              </th>
              <th>
                Отчество
                <select v-model="filters.middle_name" class="column-filter">
                  <option value="">Все</option>
                  <option v-for="name in uniqueValues('middle_name')" :key="name">{{ name }}</option>
                </select>
              </th>
              <th>
                Отдел
                <select v-model="filters.department" class="column-filter">
                  <option value="">Все</option>
                  <option v-for="dep in uniqueValues('department')" :key="dep">{{ dep }}</option>
                </select>
              </th>
              <th>
                Должность
                <select v-model="filters.position" class="column-filter">
                  <option value="">Все</option>
                  <option v-for="pos in uniqueValues('position')" :key="pos">{{ pos }}</option>
                </select>
              </th>
              <th>
                Удалённо
                <select v-model="filters.is_remote" class="column-filter">
                  <option value="">Все</option>
                  <option :value="true">Да</option>
                  <option :value="false">Нет</option>
                </select>
              </th>
              <th>
                Дата приёма
                <select v-model="filters.hire_date" class="column-filter">
                  <option value="">Все</option>
                  <option v-for="date in uniqueValues('hire_date')" :key="date">{{ date }}</option>
                </select>
              </th>
              <th>
                Зарплата
                <select v-model="filters.salary" class="column-filter">
                  <option value="">Все</option>
                  <option v-for="salary in uniqueValues('salary')" :key="salary">{{ salary }}</option>
                </select>
              </th>
              <th>Действия</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="emp in filteredEmployees" :key="emp.id">
              <td>{{ emp.id }}</td>
              <td class="fio" @click="openModal(emp.id)">{{ emp.last_name }}</td>
              <td class="fio" @click="openModal(emp.id)">{{ emp.first_name }}</td>
              <td class="fio" @click="openModal(emp.id)">{{ emp.middle_name }}</td>
              <td>{{ emp.department }}</td>
              <td>{{ emp.position }}</td>
              <td>{{ emp.is_remote ? 'Да' : 'Нет' }}</td>
              <td>{{ formatDate(emp.hire_date) }}</td>
              <td>{{ emp.salary }} ₽</td>
              <td class="actions">
                <button @click="remove(emp.id)">🗑️</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="loading" class="loading">Загрузка…</div>
        <div v-if="error" class="error">{{ error }}</div>
      </section>

      <!-- Модальное окно -->
      <div v-if="modalOpen" class="modal-overlay">
        <div class="modal">
          <h2>{{ modalEmployee.id ? 'Редактировать сотрудника' : 'Добавить сотрудника' }}</h2>
          <div class="modal-body">
            <label>Фамилия
              <input v-model="modalEmployee.last_name" placeholder="Фамилия" />
            </label>
            <label>Имя
              <input v-model="modalEmployee.first_name" placeholder="Имя" />
            </label>
            <label>Отчество
              <input v-model="modalEmployee.middle_name" placeholder="Отчество" />
            </label>
            <label>Отдел
              <select v-model="modalEmployee.department">
                <option v-for="dep in uniqueValues('department')" :key="dep">{{ dep }}</option>
              </select>
            </label>
            <label>Должность
              <select v-model="modalEmployee.position">
                <option v-for="pos in uniqueValues('position')" :key="pos">{{ pos }}</option>
              </select>
            </label>
            <label>Удалённо
              <select v-model="modalEmployee.is_remote">
                <option :value="true">Да</option>
                <option :value="false">Нет</option>
              </select>
            </label>
            <label>Дата рождения
              <input type="date" v-model="modalEmployee.birth_date" placeholder="Дата рождения" />
            </label>
            <label>Дата приёма
              <input type="date" v-model="modalEmployee.hire_date" placeholder="Дата приёма" />
            </label>
            <label>Зарплата
              <input type="number" v-model="modalEmployee.salary" placeholder="Зарплата" />
            </label>
          </div>

          <div class="modal-actions">
            <button @click="saveEmployee">Сохранить</button>
            <button @click="modalOpen = false">Отмена</button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import api from '../axios'

const employees = ref([])
const loading = ref(false)
const error = ref('')
const search = ref('')

const filters = ref({
  last_name: '',
  first_name: '',
  middle_name: '',
  department: '',
  position: '',
  is_remote: '',
  hire_date: '',
  salary: ''
})

const modalOpen = ref(false)
const modalEmployee = ref({})

async function fetchEmployees() {
  loading.value = true
  try {
    const res = await api.get('/api/employees')
    employees.value = res.data.map(emp => ({
      ...emp,
      hire_date: emp.hire_date ? new Date(emp.hire_date).toISOString().split('T')[0] : '',
      birth_date: emp.birth_date ? new Date(emp.birth_date).toISOString().split('T')[0] : ''
    }))
  } catch (e) {
    error.value = 'Ошибка загрузки сотрудников'
  } finally {
    loading.value = false
  }
}

function remove(id) {
  if (!confirm('Удалить сотрудника?')) return
  api.delete(`/api/employees/${id}`).then(() => {
    employees.value = employees.value.filter(e => e.id !== id)
  }).catch(() => alert('Ошибка удаления'))
}

function formatDate(date) {
  return date ? new Date(date).toLocaleDateString() : '—'
}

const filteredEmployees = computed(() => {
  return employees.value.filter(emp => {
    const term = search.value.toLowerCase()
    if (!(`${emp.last_name} ${emp.first_name} ${emp.middle_name}`.toLowerCase().includes(term))) return false

    for (const key in filters.value) {
      if (filters.value[key] !== '' && emp[key] != filters.value[key]) return false
    }
    return true
  })
})

function uniqueValues(field) {
  return [...new Set(employees.value.map(e => e[field]).filter(Boolean))]
}

function openModal(id) {
  const emp = employees.value.find(e => e.id === id)
  modalEmployee.value = { ...emp }
  modalOpen.value = true
}

function openAddModal() {
  modalEmployee.value = {
    id: null,
    last_name: '',
    first_name: '',
    middle_name: '',
    department: uniqueValues('department')[0] || '',
    position: uniqueValues('position')[0] || '',
    is_remote: false,
    hire_date: '',
    birth_date: '',
    salary: 0
  }
  modalOpen.value = true
}

async function saveEmployee() {
  try {
    const payload = {
      last_name: modalEmployee.value.last_name,
      first_name: modalEmployee.value.first_name,
      middle_name: modalEmployee.value.middle_name,
      department_id: getDepartmentId(modalEmployee.value.department),
      position_id: getPositionId(modalEmployee.value.position),
      is_remote: modalEmployee.value.is_remote,
      hire_date: modalEmployee.value.hire_date,
      birth_date: modalEmployee.value.birth_date,
      salary: modalEmployee.value.salary
    }

    if (modalEmployee.value.id) {
      await api.put(`/api/employees/${modalEmployee.value.id}`, payload)
    } else {
      await api.post('/api/employees', payload)
    }

    await fetchEmployees()
    modalOpen.value = false
  } catch {
    alert('Ошибка сохранения')
  }
}

function getDepartmentId(name) {
  const emp = employees.value.find(e => e.department === name)
  return emp?.department_id || 1
}
function getPositionId(name) {
  const emp = employees.value.find(e => e.position === name)
  return emp?.position_id || 1
}

onMounted(fetchEmployees)
</script>

<style scoped>
.page { display:flex; min-height:100vh; font-family: 'Nunito', sans-serif; }
.content { flex:1; padding: 36px; background: #f9fafb; box-sizing:border-box; }

.page-header h1 { margin: 0 0 8px; font-size: 24px; }
.subtitle { margin:0 0 18px; color:#64748B; }

.table-header {
  display:flex; justify-content: space-between; align-items: center; margin-bottom: 12px;
}

.add-btn {
  background-color: #4F46E5;
  color: white;
  padding: 8px 14px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
}
.add-btn:hover { background-color: #4338CA; }

.employee-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  background: #fff;
}
.employee-table th, .employee-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #e5e7eb;
  font-size: 14px;
}
.employee-table th { background: #f3f4f6; cursor: pointer; }
.employee-table th:hover { background: #e5e7eb; }

.fio { color: indigo; cursor: pointer; }
.actions button { background:none; border:none; cursor:pointer; font-size:16px; }
.actions button:hover { color:#ef4444; }

.search-input {
  width: 300px; padding: 10px 14px; border-radius:8px;
  border:1px solid #cbd5e1; background-color:#ffffff; color:#1f2937;
  font-size:14px; box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}
.search-input::placeholder { color: #9ca3af; }
.search-input:focus { outline:none; border-color:#4F46E5; box-shadow:0 0 0 2px rgba(79,70,229,0.2); }

.column-filter { width: 100%; margin-top:4px; padding:4px; border-radius:6px; border:1px solid #cbd5e1; background:#ffffff; color:#1f2937; font-size:12px; }

.modal-overlay {
  position:fixed; top:0; left:0; right:0; bottom:0; background:rgba(0,0,0,0.4); display:flex; align-items:center; justify-content:center;
}
.modal {
  background: white;
  padding: 20px;
  border-radius: 12px;
  width: 500px;
  max-width: 95%;
  box-sizing: border-box;
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 12px; 
}

.modal-body label {
  display: flex;
  flex-direction: column;
  font-size: 14px;
}

.modal-body input,
.modal-body select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
  color: #1f2937; 
  background: #fff; 
}

.modal-body input::placeholder,
.modal-body select::placeholder {
  color: #ffffff; 
  opacity: 1;
}

.modal-actions { display:flex; justify-content:flex-end; gap:12px; margin-top:16px; }
.modal-actions button { padding:8px 14px; border:none; border-radius:6px; cursor:pointer; }
.modal-actions button:first-child { background:#4F46E5; color:white; }
.modal-actions button:first-child:hover { background:#4338CA; }
.modal-actions button:last-child { background:#f3f4f6; }
.modal-actions button:last-child:hover { background:#e5e7eb; }

.loading { margin-top: 12px; color: #64748B; }
.error { margin-top: 12px; color: #EF4444; }
</style>
