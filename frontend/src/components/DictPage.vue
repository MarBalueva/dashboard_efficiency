<template>
  <div class="page">
    <Sidebar />

    <main class="content dict-layout">
      <!-- Список справочников -->
      <aside class="dict-list">
        <h3>Справочники</h3>
        <ul>
          <li
            v-for="d in dictionaries"
            :key="d.key"
            :class="{ active: currentDict.key === d.key }"
            @click="selectDict(d)"
          >
            {{ d.title }}
          </li>
        </ul>
      </aside>

      <!-- Значения справочника -->
      <section class="dict-content">
        <h2>{{ currentDict.title }}</h2>

        <button class="btn-indigo" @click="openCreate"> Добавить </button>

        <table class="dict-table">
          <thead>
            <tr>
              <th>Название</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in items" :key="item.id">
              <td>{{ item.Name }}</td>
              <td class="actions">
                <button @click="openEdit(item)">✏️</button>
                <button @click="remove(item.id)">🗑️</button>
              </td>
            </tr>
          </tbody>
        </table>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import api from '../axios'

const dictionaries = [
  { key: 'departments', title: 'Отделы', url: '/api/dict/departments' },
  { key: 'positions', title: 'Должности', url: '/api/dict/positions' },
  { key: 'access-groups', title: 'Группы доступа', url: '/api/dict/access-groups' }
]

const currentDict = ref(dictionaries[0])
const items = ref([])

async function load() {
  const res = await api.get(currentDict.value.url)
  items.value = res.data
}

function selectDict(d) {
  currentDict.value = d
  load()
}

async function openCreate() {
  const name = prompt('Название')
  const code = prompt('Код')
  if (!name) return

  await api.post(currentDict.value.url, { name, code })
  load()
}

function openEdit(item) {
  const name = prompt('Новое название', item.name)
  const code = prompt('Новый код', item.code)
  if (!name) return

  api.put(`${currentDict.value.url}/${item.id}`, { name, code })
    .then(load)
}

function remove(id) {
  if (!confirm('Удалить запись?')) return
  api.delete(`${currentDict.value.url}/${id}`)
    .then(load)
}

onMounted(load)
</script>

<style scoped>
.dict-layout {
  display: flex;
  justify-content: flex-start;
  gap: 24px;
}

.dict-list {
  width: 220px;
  background: #fff;
  border-radius: 12px;
  padding: 16px;
}

.dict-list ul {
  list-style: none;
  padding: 0;
}

.dict-list li {
  padding: 10px;
  cursor: pointer;
  border-radius: 6px;
}

.dict-list li.active,
.dict-list li:hover {
  background: #EEF2FF;
  font-weight: 600;
}

.dict-content {
  flex: 1;
  background: #fff;
  border-radius: 12px;
  padding: 24px;
}

.dict-table {
  width: 1000px;
  margin-top: 16px;
  border-collapse: collapse;
}

.dict-table th,
.dict-table td {
  padding: 10px;
  border-bottom: 1px solid #e5e7eb;
}

.actions button {
  margin-right: 6px;
  cursor: pointer;
  background: #fff;
}
</style>
