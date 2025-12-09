<template>
  <div class="page">
    <Sidebar />

    <main class="content">
      <!-- Заголовок + фильтры -->
      <header class="dashboard-header">
        <div class="title-block">
          <h1>Эффективность сотрудников</h1>
        </div>

        <div class="filters">
          <div class="filter">
            <label>Дата начала</label>
            <input type="date" v-model="dateFrom" />
          </div>

          <div class="filter">
            <label>Дата окончания</label>
            <input type="date" v-model="dateTo" />
          </div>
        </div>
      </header>

      <!-- Статистика -->
      <section class="dashboard-section">
        <div class="stats-grid">
          <div class="stat-card">
            <div class="icon-block">📘</div>
            <div class="value">{{ summary.avgLoad }} ч/мес</div>
            <div class="label">Средняя нагрузка</div>
          </div>

          <div class="stat-card">
            <div class="icon-block">⚡</div>
            <div class="value">{{ summary.avgProductivity }}</div>
            <div class="label">Средняя продуктивность</div>
          </div>

          <div class="stat-card">
            <div class="icon-block">⏱</div>
            <div class="value">{{ summary.overtime }} ч/мес</div>
            <div class="label">Переработка</div>
          </div>

          <div class="stat-card">
            <div class="icon-block">⭐</div>
            <div class="value">{{ summary.satisfaction }}</div>
            <div class="label">Удовлетворённость</div>
          </div>
        </div>
      </section>

      <!-- Графики -->
      <section class="charts-row">

        <!-- Гистограмма по отделам -->
        <div v-if="deptChartData.value?.labels?.length">
            <pre>{{ deptChartData.value }}</pre> <!-- для проверки данных -->
            <BarChartVue
                :chart-data="deptChartData.value"
                :chart-options="chartOptions"
            />
        </div>
        <!-- Топ переработок -->
        <div class="chart-card">
          <h2>Топ переработок</h2>
          <div v-for="emp in summary.topOvertime" :key="emp.name" class="overtime-card">
            <div class="emp-name">{{ emp.name }}</div>
            <div class="emp-info">
              <div>Удовлетворённость: <b>{{ emp.job_satisfaction }}</b></div>
              <div>Задачи в день: <b>{{ emp.tasks_completed_per_day }}</b></div>
              <div>Переработка: <b>{{ emp.overtime_hours }} ч/мес</b></div>
            </div>
          </div>
        </div>

      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, shallowRef, watch } from 'vue'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'
import { Bar } from 'vue-chartjs'
import api from '../axios'
import Sidebar from '../components/Sidebar.vue'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

const BarChartVue = Bar

// Даты фильтров
const dateFrom = ref('')
const dateTo = ref('')

// Сводка
const summary = ref({
  avgLoad: 0,
  avgProductivity: 0,
  overtime: 0,
  satisfaction: 0,
  deptEfficiency: [],
  topOvertime: []
})

// Chart.js данные (shallowRef чтобы избежать ошибок с реактивностью)
const deptChartData = shallowRef({
  labels: [],
  datasets: []
})

const jobLevels = ['Junior', 'Mid', 'Senior', 'Manager']
const colors = {
  Junior: '#4F46E5',
  Mid: '#0EA5E9',
  Senior: '#10B981',
  Manager: '#F59E0B'
}

const chartOptions = {
  responsive: true,
  plugins: { legend: { position: 'top' } }
}

// Загрузка данных и обновление графика
async function fetchSummary() {
  try {
    const res = await api.get('/api/dashboard/summary', {
      params: { date_from: dateFrom.value, date_to: dateTo.value }
    })

    summary.value = {
      ...res.data,
      deptEfficiency: res.data.deptEfficiency || [],
      topOvertime: res.data.topOvertime || []
    }

    const depts = summary.value.deptEfficiency || []
    if (!depts.length) {
      deptChartData.value = { labels: [], datasets: [] }
      return
    }

    const labels = [...new Set(depts.map(d => d.department))]

    const datasets = jobLevels.map(lvl => ({
      label: lvl,
      backgroundColor: colors[lvl],
      data: labels.map(dep => {
        const row = depts.find(r => r.department === dep && r.job_level === lvl)
        return row?.avg_productivity ?? 0
      })
    }))

    deptChartData.value = { labels, datasets }
  } catch (err) {
    console.error('Ошибка при загрузке сводки:', err)
    deptChartData.value = { labels: [], datasets: [] }
  }
}

// Автообновление при смене дат
watch([dateFrom, dateTo], fetchSummary)
fetchSummary()
</script>

<style scoped>
.page { display: flex; min-height: 100vh; background: #F8FAFC; }
.content { padding: 36px; flex: 1; }
.dashboard-header { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 24px; }
.title-block h1 { font-size: 26px; font-weight: 700; margin: 0; }
.filters { display: flex; gap: 20px; }
.filter { display: flex; flex-direction: column; }
.filter label { font-size: 13px; color: #64748B; margin-bottom: 4px; }
.filter input { padding: 8px 10px; border: 1px solid #cbd5e1; border-radius: 8px; background: #fff; }
.dashboard-section { display: flex; justify-content: center; padding: 20px; }
.stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 30px; }
.stat-card { background: #fff; padding: 24px; border-radius: 12px; box-shadow: 0 4px 14px rgba(0,0,0,0.06); display: flex; flex-direction: column; align-items: flex-start; }
.icon-block { font-size: 28px; margin-bottom: 10px; }
.value { font-size: 28px; font-weight: 700; margin-bottom: 6px; color: #0f172a; }
.label { font-size: 14px; color: #64748b; }
.charts-row { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-top: 40px; }
.chart-card { background: #fff; padding: 24px; border-radius: 12px; box-shadow: 0 4px 14px rgba(0,0,0,0.06); height: auto; }
.chart-card h2 { margin-bottom: 16px; font-size: 18px; font-weight: 700; }
.overtime-card { padding: 16px; border-radius: 10px; background: #F8FAFC; margin-bottom: 14px; display: flex; justify-content: space-between; align-items: center; }
.emp-name { font-size: 16px; font-weight: 600; }
.emp-info div { font-size: 14px; color: #475569; }
</style>
