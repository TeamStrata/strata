<script setup>
import { ref, reactive } from 'vue'
import { watch } from 'vue'
import { computed, defineAsyncComponent } from 'vue'
import { apiFetch } from '@/api/request'

const chartTypes = ['Line', 'Area', 'Column', 'Bar', 'Scatter']
const selectedChart = ref(chartTypes[0])
const isCreatingNewChart = ref(false)

const savedQueries = ref([])
const columns = ref([])

const chartTitle = ref('')
const savedChartTitles = ref([])
const selectedChartTitle = ref('')

const seriesSections = reactive([])
const series = ref([
  {
    query: {},
    x_column: '',
    y_column: ''
  }
])
const chartComponent = computed(() => {
  const map = {
    Line: defineAsyncComponent(() => import('./charts/LineChart.vue')),
    Area: defineAsyncComponent(() => import('./charts/AreaChart.vue')),
    Column: defineAsyncComponent(() => import('./charts/ColumnChart.vue')),
    Bar: defineAsyncComponent(() => import('./charts/BarChart.vue')),
    Scatter: defineAsyncComponent(() => import('./charts/ScatterChart.vue')),
  }

  return map[selectedChart.value] || null
})
// Fetch saved queries
apiFetch('/queries')
  .then(async (response) => {
    if (!response.ok) {
      toasRef.value?.showToast(
        'There was an error when loading saved queries',
        ToastTypes.FAIL
      )
      throw new Error('Error loading saved queries')
    }

    savedQueries.value = await response.json()

    // Only continue if queries exist
    if (savedQueries.value.length > 0) {
      series.value[0].query = savedQueries.value[0]

      // Now fetch columns for selected query
      return apiFetch(`/query/${series.value[0].query.id}/execute`)
    } else {
      throw new Error('No queries found')
    }
  })
 .then(async (response) => {
    if (!response.ok) {
      toasRef.value?.showToast(
        'There was an error when executing the query',
        ToastTypes.FAIL
      )
      throw new Error('Error executing query')
    }

    const data = await response.json()
    columns.value = Object.keys(data[0] || {})
    series.value[0].x_column = columns.value[0]
    series.value[0].y_column = columns.value[1]

    return apiFetch('/charts/titles')
  })
  .then(async (response) => {
    if (!response.ok) throw new Error('Failed to fetch chart titles')
    savedChartTitles.value = await response.json()
  })
  .catch((error) => {
    console.error('Initialization error:', error)
  })
const addSeriesSection = () => {
  seriesSections.push({
    query: savedQueries.value[0] || {},
    xColumn: columns.value[0] || '',
    yColumn: columns.value[1] || '',
     chartData: [],
  })
}

const removeChartSection = (index) => {
  seriesSections.splice(index, 1)
}

const chartData = ref([])

const generateChart = async (index) => {
  const section = seriesSections[index]
  if (!section.query?.id) return

  try {
    const response = await apiFetch(`/query/${section.query.id}/execute`)
    if (!response.ok) throw new Error('Failed to execute query')

    const data = await response.json()
    seriesSections[index].chartData = data

    console.log("Generating chart with:", {
      chartType: selectedChart.value,
      x: section.xColumn,
      y: section.yColumn,
      data
    })
  } catch (err) {
    console.error('Error generating chart:', err)
  }
}

const saveChart = async () => {
  const payload = {
    chart_title: chartTitle.value,
    chart_type: selectedChart.value.toLowerCase()|| 'line',
  }

  try {
    const response = await apiFetch('/chart', 'POST', JSON.stringify(payload), "application/sql")
    if (!response.ok) throw new Error('Failed to save chart')

    alert('Chart saved successfully')
    fetchSavedChartTitles()
  } catch (err) {
    console.error('Error saving chart:', err)
  }
}

const loadChartFromDB = async () => {
  if (!selectedChartTitle.value) return

  try {
    const response = await apiFetch(`/chart/${selectedChartTitle.value}`)
    if (!response.ok) throw new Error('Failed to load chart')
    const data = await response.json()

    chartTitle.value = data.title
    selectedChart.value = data.chartType
    seriesSections.splice(0, seriesSections.length, ...data.series)
  } catch (err) {
    console.error('Error loading chart:', err)
  }
}

watch(seriesSections, (sections) => {
  sections.forEach((section, index) => {
    watch(
      () => section.query,
      async (newQuery) => {
        if (!newQuery || !newQuery.id) return

        try {
          const response = await apiFetch(`/query/${newQuery.id}/execute`)
          if (!response.ok) throw new Error('Failed to execute query')

          const data = await response.json()
          const cols = Object.keys(data[0] || {})

          columns.value = cols

          if (!section.xColumn) section.xColumn = cols[0] || ''
          if (!section.yColumn) section.yColumn = cols[1] || ''
        } catch (err) {
          console.error(`Error loading columns for query ${newQuery.id}:`, err)
        }
      },
      { immediate: true }
    )
  })
}, { deep: true })

const fetchSavedChartTitles = async () => {
  try {
    const response = await apiFetch('/charts/titles')
    if (!response.ok) throw new Error('Failed to fetch chart titles')
    savedChartTitles.value = await response.json()
  } catch (err) {
    console.error('Error fetching titles:', err)
  }
}
</script>


<template>
  <div class="max-w-full mx-auto mt-10 p-6 bg-white shadow-md rounded-lg">
    <h1 class="text-2xl font-bold mb-6 text-gray-800">Chart Viewer</h1>

    <!-- Chart Type Selection -->
    <div class="mb-4">
      <label class="text-gray-600 font-medium">Chart Type:</label>
      <select
        v-model="selectedChart"
        class="w-full p-2 border border-gray-300 rounded"
      >
        <option v-for="type in chartTypes" :key="type">{{ type }}</option>
      </select>
    </div>

    <!-- New or Load Chart Options -->
    <div class="flex flex-wrap gap-4 mb-6">
      <div class="flex flex-col">
        <label class="text-gray-600 font-medium mb-1">Load Saved Chart:</label>
        <select
          v-model="selectedChartTitle"
          @change="loadChartFromDB"
          class="p-2 border border-gray-300 rounded"
        >
          <option value="">Select a chart</option>
          <option v-for="title in savedChartTitles" :key="title" :value="title">
            {{ title }}
          </option>
        </select>
      </div>
      <button
        @click="isCreatingNewChart = true"
        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
      >
        Create New Chart
      </button>
    </div>

    <!-- Only show this if user is creating a new chart -->
    <div v-if="isCreatingNewChart">
      <!-- Chart Title Input -->
      <div class="mb-4">
        <label class="text-gray-600 font-medium">Chart Title:</label>
        <input
          v-model="chartTitle"
          class="w-full p-2 border border-gray-300 rounded"
          placeholder="Enter chart title"
        />
      </div>

      <!-- Series Sections (Horizontally Scrollable) -->
      <div class="flex gap-4 overflow-x-auto pb-4">
        <div
          v-for="(section, index) in seriesSections"
          :key="index"
          class="min-w-[300px] shrink-0 border border-gray-300 p-4 rounded-lg bg-gray-50 relative"
        >
          <h2 class="text-md font-semibold text-gray-700 mb-2">Series {{ index + 1 }}</h2>

          <!-- Remove Button -->
          <button
            @click="removeChartSection(index)"
            class="absolute top-2 right-2 text-red-500 hover:text-red-700 text-sm"
          >
            ✕
          </button>

          <!-- Query -->
          <div class="mb-2">
            <label class="block text-gray-600 text-sm mb-1">Saved Query:</label>
            <select
              v-model="section.query"
              class="w-full p-1 border border-gray-300 rounded text-sm"
            >
              <option v-for="query in savedQueries" :key="query.id" :value="query">
                {{ query.name || query.id }}
              </option>
            </select>
          </div>

          <!-- X Column -->
          <div class="mb-2">
            <label class="block text-gray-600 text-sm mb-1">X Column:</label>
            <select
              v-model="section.xColumn"
              class="w-full p-1 border border-gray-300 rounded text-sm"
            >
              <option v-for="col in columns" :key="col">{{ col }}</option>
            </select>
          </div>

          <!-- Y Column -->
          <div class="mb-2">
            <label class="block text-gray-600 text-sm mb-1">Y Column:</label>
            <select
              v-model="section.yColumn"
              class="w-full p-1 border border-gray-300 rounded text-sm"
            >
              <option v-for="col in columns" :key="col">{{ col }}</option>
            </select>
          </div>

          <!-- Generate Button -->
          <button
            @click="generateChart(index)"
            class="mt-2 w-full bg-blue-500 hover:bg-blue-600 text-white text-sm py-1 rounded"
          >
            Generate
          </button>
        </div>
      </div>

      <!-- Add & Save Buttons -->
      <div class="mt-6 flex flex-wrap gap-4">
        <button
          @click="addSeriesSection"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
        >
          Add Series
        </button>

        <button
          @click="saveChart"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        >
          Save Chart
        </button>
      </div>
      <!-- Render Chart -->
      
      <div class="mt-6" v-if="chartComponent && seriesSections.length">
        <div v-for="(section, index) in seriesSections" :key="index">
          <component
          :is="chartComponent"
          :series="seriesSections.filter(s => s.chartData?.length)"
          />
        </div>
    </div>
    </div>
  </div>
</template>
