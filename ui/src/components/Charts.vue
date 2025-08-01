<script setup>
import { ref, reactive } from 'vue'
import { watch } from 'vue'
import { computed, defineAsyncComponent } from 'vue'
import { apiFetch } from '@/api/request'
import Input from './ui/input/Input.vue';
import Label from './ui/label/Label.vue';
import { File, FolderOpen, Plus, Save, X } from 'lucide-vue-next';

import Dialog from './ui/dialog/Dialog.vue';
import DialogContent from './ui/dialog/DialogContent.vue';
import DialogHeader from './ui/dialog/DialogHeader.vue';
import DialogTitle from './ui/dialog/DialogTitle.vue';

import Toast, { ToastTypes } from "./Toast.vue";
import Separator from './ui/separator/Separator.vue';

const toastRef = ref(null);
const chartTypes = ['Line', 'Area', 'Column', 'Bar', 'Scatter']
const selectedChart = ref(chartTypes[0])
const isCreatingNewChart = ref(false)

const savedQueries = ref([])
const columns = ref([])

const chartTitle = ref('')
const savedChartTitles = ref([])
const selectedChartTitle = ref('')
const isChartLoaded = ref(false)

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
      toastRef.value?.showToast(
        'There was an error when loading saved queries',
        ToastTypes.FAIL
      )
      throw new Error('Error loading saved queries')
    }

    savedQueries.value = await response.json()

    if (savedQueries.value.length > 0) {
      series.value[0].query = savedQueries.value[0]

      return apiFetch(`/query/${series.value[0].query.id}/execute`)
    } else {
      throw new Error('No queries found')
    }
  })
  .then(async (response) => {
    if (!response.ok) {
      toastRef.value?.showToast(
        'There was an error when executing the query',
        ToastTypes.FAIL
      )
      throw new Error('Error executing query')
    }

    const data = await response.json()
    columns.value = Object.keys(data[0] || {})
    series.value[0].x_column = columns.value[0]
    series.value[0].y_column = columns.value[1]

    return apiFetch('/charts')
  })
  .then(async (response) => {
    if (!response.ok) throw new Error('Failed to fetch chart titles')
    savedChartTitles.value = await response.json()
  })
  .catch((error) => {
    console.error('Initialization error:', error)
  })

const addSeriesSection = () => {
  const section = reactive({
    query: savedQueries.value[0] || {},
    xColumn: '',
    yColumn: '',
    chartData: [],
    columns: [],
  })
  watch(
    () => section.query,
    async (newQuery) => {
      if (!newQuery || !newQuery.id) return

      try {
        const response = await apiFetch(`/query/${newQuery.id}/execute`)
        if (!response.ok) throw new Error('Failed to execute query')

        const data = await response.json()
        const cols = Object.keys(data[0] || {})

        section.columns = cols
        if (!section.xColumn) section.xColumn = cols[0] || ''
        if (!section.yColumn) section.yColumn = cols[1] || ''
      } catch (err) {
        console.error(`Error loading columns for query ${newQuery.id}:`, err)
      }
    },
    { immediate: true }
  )

  seriesSections.push(section)
}


const removeChartSection = (index) => {
  let spliced = seriesSections.splice(index, 1)
  spliced.forEach(s => {
    apiFetch(`/chart/${selectedChartTitle.value.id}/series/${s.id}`, 'DELETE')
      .then(response => {
        if (!response.ok) throw new Error('Failed to delete chart section')
        console.log(`Section ${s.id} removed successfully`)
      })
      .catch(err => {
        console.error(`Error removing section ${s.id}:`, err)
      })
  })
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
    title: chartTitle.value,
    type: selectedChart.value.toLowerCase(),
  }

  try {
    if (!payload.title) {
      toastRef.value?.showToast('Missing chart title', ToastTypes.FAIL)
      throw new Error('Missing chart title')
    }

    if (seriesSections.length < 1) {
      toastRef.value?.showToast(
        'At least one series is required',
        ToastTypes.FAIL
      )
      throw new Error('At least one series is required')
    }

    const response = await apiFetch('/chart', 'POST', JSON.stringify(payload), 'application/sql')
    if (!response.ok) throw new Error('Failed to save chart')

    const data = await response.json()
    const chartId = data

    if (!chartId) {
      throw new Error('No chart ID returned from server')
    }

    const seriesPayload = seriesSections.map(section => ({
      id: section.id,
      chart_id: chartId,
      query_id: section.query?.id,
      x_col_name: section.xColumn,
      y_col_name: section.yColumn,
    }))

    console.log("Posting chartSeries links:");
    seriesPayload.forEach(async c => {
      console.log(c);
      const seriesResponse = await apiFetch(`/chart/${c.chart_id}/series`, 'POST', JSON.stringify(c), 'application/json')
      if (!seriesResponse.ok) throw new Error('Failed to save series')
      else
        alert('Chart and series saved successfully')
    })



    fetchSavedChartTitles()
  } catch (err) {
    console.error('Error saving chart and series:', err)
  }
}
function capitalizeFirstLetter(str) {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()
}
const loadChartFromDB = async () => {
  if (!selectedChartTitle.value || !selectedChartTitle.value.id) return;

  try {
    // Fetch the chart info
    const chartRes = await apiFetch(`/chart/${selectedChartTitle.value.id}`);
    if (!chartRes.ok) throw new Error('Failed to load chart');
    const chartData = await chartRes.json();

    chartTitle.value = chartData.title;
    selectedChart.value = capitalizeFirstLetter(chartData.type);

    // Fetch the series separately (assuming /chart/:id/series endpoint)
    const seriesRes = await apiFetch(`/chart/${selectedChartTitle.value.id}/series`);
    if (!seriesRes.ok) throw new Error('Failed to load chart series');
    const seriesList = await seriesRes.json();

    // Load seriesSections with loaded series + queries + chartData
    const loadedSeries = await Promise.all(
      seriesList.map(async (seriesItem) => {
        const query = savedQueries.value.find(q => q.id === seriesItem.query_id);
        if (!query) return null;

        // Execute query to get chart data for this series
        const result = await apiFetch(`/query/${query.id}/execute`);
        const chartData = await result.json();

        return {
          id: seriesItem.id,
          query,
          xColumn: seriesItem.x_col_name,
          yColumn: seriesItem.y_col_name,
          chart_id: seriesItem.chart_id,
          chartData,
          columns: Object.keys(chartData[0] || {}), // dynamically get columns for selects
        };
      })
    );
    isChartLoaded.value = true
    isCreatingNewChart.value = false
    seriesSections.splice(0, seriesSections.length, ...loadedSeries.filter(Boolean));
    isCreatingNewChart.value = false;

  } catch (err) {
    console.error('Error loading chart:', err);
  }
};

const fetchSavedChartTitles = async () => {
  try {
    const response = await apiFetch('/charts')
    if (!response.ok) throw new Error('Failed to fetch chart titles')
    savedChartTitles.value = await response.json()
  } catch (err) {
    console.error('Error fetching titles:', err)
  }
}

async function deleteChart(cid) {
  apiFetch(`/chart/${cid}`, 'DELETE').then(res => {
    if (res.ok) {
      toastRef.value?.showToast('Chart deleted successfully', ToastTypes.SUCCESS)
      fetchSavedChartTitles()
    } else {
      toastRef.value?.showToast('Failed to delete chart', ToastTypes.FAIL)
    }
  })
}

const isLoadOpen = ref(false);
</script>


<template>
  <Toast ref="toastRef" />

  <!-- load dialog -->
  <Dialog v-model:open="isLoadOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Load Chart</DialogTitle>
      </DialogHeader>
      <div>
        <div class="flex items-center justify-between font-semibold text-gray-700 px-2 py-1 border-b mb-2">
          <span class="w-1/3">Title</span>
          <span class="w-1/3 text-right">Type</span>
          <span class="w-8"></span>
        </div>
        <div class="flex flex-col">
          <div class="flex items-center justify-between group hover:bg-muted p-2 cursor-pointer"
            v-for="chart in savedChartTitles"
            @click="selectedChartTitle = chart; loadChartFromDB(); isLoadOpen = false">
            <p class="w-1/3 truncate">{{ chart.title }}</p>
            <p class="w-1/3 text-right">{{ chart.type }}</p>
            <div class="w-8 invisible group-hover:visible flex justify-end" @click="deleteChart(chart.id)">
              <X></X>
            </div>
          </div>
        </div>
      </div>
    </DialogContent>
  </Dialog>


  <div class="flex h-full w-full">
    <!-- left control panel -->
    <div class="w-[350px] h-[calc(100vh-70px)] flex flex-col py-2 px-4 gap-4 min-h-0">
      <!-- save/load controls -->
      <div class="flex gap-1 justify-end">
        <div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm">
          <File></File>
        </div>
        <div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm" @click="isLoadOpen = true">
          <FolderOpen></FolderOpen>
        </div>
        <div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm" @click="saveChart">
          <Save :size="24" class="hover:cursor-pointer"></Save>
        </div>
      </div>

      <!-- name and stuff -->
      <div>
        <Label class="mb-1">Chart Title</Label>
        <Input v-model="chartTitle" placeholder="Enter chart title" />
      </div>

      <!-- type selection -->
      <div>
        <Label class="mb-1">Chart Type:</Label>
        <select v-model="selectedChart" class="w-full p-2 border border-gray-300 rounded">
          <option v-for="type in chartTypes" :key="type">{{ type }}</option>
        </select>
      </div>

      <!-- series management -->
      <div class="flex justify-between items-center">
        <p class="leading-7 [&:not(:first-child)]:mt-6">Series</p>
        <Plus @click="addSeriesSection" class="cursor-pointer"></Plus>
      </div>

      <!-- series scrollable container -->
      <div class="flex-1 overflow-y-scroll gap-4 flex flex-col">
        <div v-for="(section, index) in seriesSections" :key="index"
          class="min-w-[300px] shrink-0 border border-gray-300 p-4 rounded-lg bg-gray-50 relative">
          <div class="flex justify-between items-center mb-2">
            <h2 class="text-md font-semibold text-gray-700">Series {{ index + 1 }}</h2>
            <!-- Remove Button -->
            <button @click="removeChartSection(index)" class="text-red-500 hover:text-red-700 cursor-pointer text-sm">
              ✕
            </button>
          </div>

          <!-- Query -->
          <div class="mb-2">
            <label class="block text-gray-600 text-sm mb-1">Saved Query:</label>
            <select v-model="section.query" class="w-full p-1 border border-gray-300 rounded text-sm">
              <option v-for="query in savedQueries" :key="query.id" :value="query">
                {{ query.name || query.id }}
              </option>
            </select>
          </div>

          <!-- X Column -->
          <div class="mb-2">
            <label class="block text-gray-600 text-sm mb-1">X Column:</label>
            <select v-model="section.xColumn" class="w-full p-1 border border-gray-300 rounded text-sm">
              <option v-for="col in section.columns" :key="col">{{ col }}</option>
            </select>
          </div>

          <!-- Y Column -->
          <div class="mb-2">
            <label class="block text-gray-600 text-sm mb-1">Y Column:</label>
            <select v-model="section.yColumn" class="w-full p-1 border border-gray-300 rounded text-sm">
              <option v-for="col in section.columns" :key="col">{{ col }}</option>
            </select>
          </div>

          <!-- Generate Button -->
          <button @click="generateChart(index)"
            class="mt-2 w-full bg-blue-500 hover:bg-blue-600 text-white text-sm py-1 rounded">
            Generate
          </button>
        </div>

      </div>

    </div>
    <Separator orientation="vertical"></Separator>
    <!-- Render Chart -->
    <div class="flex-1 ml-4 my-auto" v-if="chartComponent && seriesSections.length">
      <div v-if="chartComponent && seriesSections.some(s => s.chartData?.length)">
        <component :is="chartComponent" :series="seriesSections.filter(s => s.chartData?.length)" />
      </div>
    </div>
  </div>
</template>
