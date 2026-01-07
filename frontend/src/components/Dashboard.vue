<template>
  <div class="page">
    <Sidebar />

    <main class="content">
        <div class="dashboard-container">
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

      <section class="dashboard-section">
        <div class="stats-grid">
          <div class="stat-card">
            <div class="icon-block">📘</div>
            <div class="value">{{ formatNumber(summary.avgLoad, 1) }} ч/мес</div>
            <div class="label">Средняя нагрузка</div>
          </div>

          <div class="stat-card">
            <div class="icon-block">⚡</div>
            <div class="value">{{ formatNumber(summary.avgProductivity, 2) }}</div>
            <div class="label">Средняя продуктивность</div>
          </div>

          <div class="stat-card">
            <div class="icon-block">⏱</div>
            <div class="value">{{ formatNumber(summary.overtime, 1) }} ч/мес</div>
            <div class="label">Переработка</div>
          </div>

          <div class="stat-card">
            <div class="icon-block">⭐</div>
            <div class="value">{{ formatNumber(summary.satisfaction, 1) }}</div>
            <div class="label">Удовлетворённость</div>
          </div>
        </div>
      </section>

      <section class="charts-row">
        <div class="chart-card">
          <h2>Эффективность по отделам</h2>

          <div v-if="chartSeries.length && chartCategories.length">
            <ApexChart
              type="bar"
              height="360"
              :options="chartOptions"
              :series="chartSeries"
            />
          </div>
          <div v-else class="empty-note">Нет данных для отображения графика</div>
        </div>

        <div class="chart-card">
          <h2>Топ переработок</h2>

          <div v-if="summary.topOvertime.length === 0" class="empty-note">Нет данных</div>

          <div v-for="emp in summary.topOvertime" :key="emp.name" class="overtime-card">
            <div class="emp-left">
              <div class="emp-name">Сотрудник: {{ emp.name }}</div>
            </div>
            <div class="emp-right">
              <div>Удовлетворённость: <b>{{ emp.job_satisfaction }}</b></div>
              <div>Задачи в день: <b>{{ emp.tasks_completed_per_day }}</b></div>
              <div>Переработка: <b>{{ emp.overtime_hours }} ч/мес</b></div>
            </div>
          </div>
        </div>

        <div class="chart-card">
        <h2>Нагрузка и переработка по месяцам</h2>
        <ApexChart
            type="line"
            height="360"
            :options="monthlyChartOptions"
            :series="monthlySeries"
            v-if="monthlySeries.length"
        />
        <div v-else class="empty-note">Нет данных</div>
        </div>

        <div class="chart-card">
        <h2>Топ сотрудников по эффективности</h2>
        <div v-if="!summary.topEfficiency || summary.topEfficiency.length === 0" class="empty-note">
            Нет данных
        </div>
        <div v-for="emp in summary.topEfficiency" :key="emp.name" class="overtime-card">
            <div class="emp-left">
            <div class="emp-name">Сотрудник: {{ emp.name }}</div>
            </div>
            <div class="emp-right">
            <div>Продуктивность: <b>{{ formatNumber(emp.productivity, 2) }}</b></div>
            <div>Задачи в день: <b>{{ emp.tasks }}</b></div>
            <div>Переработка: <b>{{ formatNumber(emp.overtime, 1) }} ч/мес</b></div>
            </div>
        </div>
        </div>

      </section>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import ApexChart from 'vue3-apexcharts'
import api from '../axios'
import Sidebar from '../components/Sidebar.vue'

// --- Даты по умолчанию: текущий месяц ---
const now = new Date()
const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
const lastDay = new Date(now.getFullYear(), now.getMonth() + 1, 0)

function formatDate(d) {
  const month = (d.getMonth() + 1).toString().padStart(2, '0')
  const day = d.getDate().toString().padStart(2, '0')
  return `${d.getFullYear()}-${month}-${day}`
}

const dateFrom = ref(formatDate(firstDay))
const dateTo = ref(formatDate(lastDay))

// --- Сводка ---
const summary = ref({
  avgLoad: 0,
  avgProductivity: 0,
  overtime: 0,
  satisfaction: 0,
  deptEfficiency: [],
  topOvertime: [],
  monthlyStats: [],
  topEfficiency: []  
})

// --- Должности динамические ---
const jobLevels = ref([])

async function fetchPositions() {
  try {
    const res = await api.get('/api/dict/positions')
    jobLevels.value = res.data.map(p => p.name) // assuming {id, name}
  } catch (err) {
    console.error('Ошибка загрузки должностей:', err)
    jobLevels.value = []
  }
}

const chartCategories = computed(() => {
  const rows = summary.value.deptEfficiency || []
  return [...new Set(rows.map(r => r.department))] 
})

const barColors = ['#FF8F6B', '#4F46E5', '#FFD66B']
function getColor(index) { return barColors[index % barColors.length] }

const chartSeries = computed(() => {
  const rows = summary.value.deptEfficiency || []
  const cats = chartCategories.value
  if (!cats.length || !rows.length) return []

  const levels = [...new Set(rows.map(r => r.job_level))] 

  return levels.map((lvl, idx) => {
    const data = cats.map(dep => {
      const row = rows.find(r => r.department === dep && r.job_level === lvl)
      return row && row.avg_productivity !== undefined ? Number(row.avg_productivity.toFixed(2)) : 0
    })
    return { name: lvl, data, color: getColor(idx) }
  }).filter(s => s.data.some(v => v > 0))
})


// --- Опции графика ---
const chartOptions = computed(() => ({
  chart: { type: 'bar', stacked: false, toolbar: { show: true } },
  plotOptions: { bar: { horizontal: false, columnWidth: '55%', borderRadius: 6, borderRadiusApplication: 'end', dataLabels: { position: 'top' } } },
  xaxis: { categories: chartCategories.value, labels: { rotate: -15 } },
  yaxis: { title: { text: 'Эффективность' } },
  legend: { position: 'top' },
  dataLabels: { enabled: false },
  tooltip: { y: { formatter: val => val } }
}))

// --- График нагрузки и переработки ---
const monthlyCategories = computed(() => (summary.value.monthlyStats || []).map(m => m.month))
const monthlySeries = computed(() => {
  const rows = summary.value.monthlyStats || []
  return [
    { name: "Нагрузка", data: rows.map(r => Number(r.load.toFixed(2))), color: "#4F46E5" },
    { name: "Переработка", data: rows.map(r => Number(r.overtime.toFixed(2))), color: "#FF8F6B" }
  ]
})
const monthlyChartOptions = computed(() => ({
  chart: { type: "line", toolbar: { show: true } },
  stroke: { curve: "smooth", width: 3 },
  markers: { size: 4, hover: { size: 7 } },
  xaxis: { categories: monthlyCategories.value },
  yaxis: { title: { text: "Значения" } },
  tooltip: { shared: true, intersect: false },
  legend: { position: "top" }
}))

// --- Получение сводки ---
async function fetchSummary() {
  try {
    const res = await api.get('/api/dashboard/summary', {
      params: { date_from: dateFrom.value, date_to: dateTo.value }
    })

    summary.value = {
      avgLoad: Number(res.data.avgLoad?.toFixed(2) ?? 0),
      avgProductivity: Number(res.data.avgProductivity?.toFixed(2) ?? 0),
      overtime: Number(res.data.overtime?.toFixed(2) ?? 0),
      satisfaction: Number(res.data.satisfaction?.toFixed(2) ?? 0),
      deptEfficiency: res.data.deptEfficiency ?? [],
      topOvertime: (res.data.topOvertime ?? []).map(emp => ({
        ...emp,
        overtime_hours: Number(emp.overtime_hours.toFixed(2))
      })),
      monthlyStats: (res.data.monthlyStats ?? []).map(m => ({
        month: m.month,
        load: Number(m.load.toFixed(2)),
        overtime: Number(m.overtime.toFixed(2))
      })),
      topEfficiency: (res.data.topEfficiency ?? []).map(emp => ({
        ...emp,
        productivity: Number(emp.productivity.toFixed(2)),
        overtime: Number(emp.overtime.toFixed(2))
      }))
    }
  } catch (err) {
    console.error('Ошибка при fetchSummary:', err)
  }
}

watch([dateFrom, dateTo], fetchSummary)

onMounted(async () => {
  await fetchPositions()
  fetchSummary()
})

function formatNumber(v, digits = 2) {
  if (v === null || v === undefined) return '-'
  return Number(v).toFixed(digits)
}
</script>

<style scoped>
/* .dashboard-container {
  background: #fff;
  border-radius: 16px;
  padding: 32px;
  width: calc(100% - 40px);
  margin-left: 20px; 
  margin-right: 20px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.05);
  min-height: calc(100vh - 72px);
} */

.page { background: #F8FAFC;
  border-radius: 16px;
  padding: 32px;
  width: 100%;
  margin-left: 20px; 
  margin-right: 20px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.05);
  min-height: calc(100vh - 72px); }
.content { flex: 1; }

/* header */
.dashboard-header { display:flex; justify-content:space-between; align-items:flex-end; margin-bottom:24px; }
.title-block h1 { font-size:26px; font-weight:700; margin:0; }

/* filters */
.filters { display:flex; gap:20px; }
.filter { display:flex; flex-direction:column; }
.filter label { font-size:13px; color:#64748B; margin-bottom:4px; }
.filter input { padding:8px 10px; border:1px solid #cbd5e1; border-radius:8px; background:#fff }
.filter input[type="date"] {
  color: #000; 
  background-color: #fff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 10px;
}

.filter input[type="date"]::-webkit-calendar-picker-indicator {
  filter: invert(0); 
}

/* stats */
.dashboard-section { display:flex; justify-content:center; padding:20px; }
.stats-grid { display:grid; grid-template-columns:repeat(4,1fr); gap:20px; margin-bottom:30px; width:100%; max-width:1100px; }
.stat-card { background:#fff; padding:24px; border-radius:12px; box-shadow:0 4px 14px rgba(0,0,0,0.06); display:flex; flex-direction:column; align-items:flex-start; }
.icon-block { font-size:28px; margin-bottom:10px; }
.value { font-size:28px; font-weight:700; margin-bottom:6px; color:#0f172a; }
.label { font-size:14px; color:#64748b; }

/* charts row */
.charts-row { display:grid; grid-template-columns:1fr 1fr; gap:24px; margin-top:40px; }
.chart-card { background:#fff; padding:24px; border-radius:12px; box-shadow:0 4px 14px rgba(0,0,0,0.06); }
.chart-card h2 { margin-bottom:16px; font-size:18px; font-weight:700; }
.empty-note { color:#64748B; padding:16px; }

/* overtime cards */
.overtime-card { padding:16px; border-radius:10px; background:#F8FAFC; margin-bottom:14px; display:flex; justify-content:space-between; align-items:center; }
.emp-name { font-size:16px; font-weight:600; }
.emp-info div { font-size:14px; color:#475569; }

/* responsive */
@media (max-width: 900px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
  .charts-row { grid-template-columns: 1fr; }
}
</style>
