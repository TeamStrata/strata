<script setup>
import { ref, reactive } from 'vue'
import { watch } from 'vue'
import { computed, defineAsyncComponent } from 'vue'
import { apiFetch } from '@/api/request'
import Input from './ui/input/Input.vue';
import Label from './ui/label/Label.vue';
import { Check, File, FolderOpen, Pencil, Plus, Save, Watch, X } from 'lucide-vue-next';

import Dialog from './ui/dialog/Dialog.vue';
import DialogContent from './ui/dialog/DialogContent.vue';
import DialogHeader from './ui/dialog/DialogHeader.vue';
import DialogTitle from './ui/dialog/DialogTitle.vue';

import Toast, { ToastTypes } from "./Toast.vue";
import Separator from './ui/separator/Separator.vue';
import Swal from 'sweetalert2';
import DialogDescription from './ui/dialog/DialogDescription.vue';

const toastRef = ref(null);
const chartTypes = ['Line', 'Area', 'Column', 'Bar', 'Scatter']
const isCreatingNewChart = ref(false)

const savedQueries = ref([])
const columns = ref([])

const savedChartTitles = ref([])
const selectedChartTitle = ref('')
const isChartLoaded = ref(false)

const series = ref([
  {
    query: {},
    x_column: '',
    y_column: ''
  }
])

const chartFull = reactive({
  chart_id: null,
  title: "",
  type: chartTypes[0],
  xColumn: '',
  yColumn: '',
  chartData: [],
})

const chartComponent = computed(() => {
  const map = {
    Line: defineAsyncComponent(() => import('./charts/LineChart.vue')),
    Area: defineAsyncComponent(() => import('./charts/AreaChart.vue')),
    Column: defineAsyncComponent(() => import('./charts/ColumnChart.vue')),
    Bar: defineAsyncComponent(() => import('./charts/BarChart.vue')),
    Scatter: defineAsyncComponent(() => import('./charts/ScatterChart.vue')),
  }

  return map[chartFull.type] || null
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
    data: [],
    columns: [],
    name: "Series " + (chartFull.chartData.length + 1)
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

  chartFull.chartData.push(section)
  generateChart(chartFull.chartData.length - 1)
}


const removeChartSection = (index) => {
  chartFull.chartData.splice(index, 1)
}

const generateChart = async (index) => {
  const section = chartFull.chartData[index]
  if (!section.query?.id) return

  try {
    const response = await apiFetch(`/query/${section.query.id}/execute`)
    if (!response.ok) throw new Error('Failed to execute query')

    const data = await response.json()
    chartFull.chartData[index].chartData = data

  } catch (err) {
    console.error('Error generating chart:', err)
  }
}

const saveChart = async () => {
  const payload = {
    title: chartFull.title,
    type: chartFull.type.toLowerCase(),
    x_axis: chartFull.xColumn,
    y_axis: chartFull.yColumn
  }

  try {
    if (!payload.title) {
      toastRef.value?.showToast('Missing chart title', ToastTypes.FAIL)
      throw new Error('Missing chart title')
    }

    if (chartFull.chartData.length < 1) {
      toastRef.value?.showToast(
        'At least one series is required',
        ToastTypes.FAIL
      )
      throw new Error('At least one series is required')
    }

    const response = await apiFetch('/chart', 'POST', JSON.stringify(payload), 'application/json')
    if (!response.ok) throw new Error('Failed to save chart')

    const data = await response.json()
    const chartId = data

    if (!chartId) {
      throw new Error('No chart ID returned from server')
    }

    const seriesPayload = chartFull.chartData.map(section => ({
      id: section.id,
      chart_id: chartId,
      query_id: section.query?.id,
      xColumn: section.xColumn,
      yColumn: section.yColumn,
      name: section.name,
    }))

    console.log("Posting chartSeries links:");
    const seriesResponse = await apiFetch(`/chart/${chartId}/series`, 'POST', JSON.stringify(seriesPayload), 'application/json')
    if (!seriesResponse.ok) throw new Error('Failed to save series')
    else
      Swal.fire({
        title: 'Success!',
        text: 'Chart and series saved successfully.',
        icon: 'success',
        timer: 2000,
        showConfirmButton: false
      });




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

    resetChart();

    chartFull.title = chartData.title;
    chartFull.type = capitalizeFirstLetter(chartData.type);
    chartFull.xColumn = chartData.x_axis || '';
    chartFull.yColumn = chartData.y_axis || '';
    chartFull.chart_id = chartData.id;

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
          name: seriesItem.name,
          id: seriesItem.id,
          query,
          xColumn: seriesItem.xColumn,
          yColumn: seriesItem.yColumn,
          chart_id: seriesItem.chart_id,
          data: chartData,
          columns: Object.keys(chartData[0] || {}), // dynamically get columns for selects
        };
      })
    );
    isChartLoaded.value = true
    isCreatingNewChart.value = false
    console.log()
    chartFull.chartData.splice(0, chartFull.chartData.length, ...loadedSeries.filter(Boolean));
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

function resetChart() {
  chartFull.title = ''
  chartFull.type = chartTypes[0]
  chartFull.xColumn = ''
  chartFull.yColumn = ''
  chartFull.chartData.splice(0, chartFull.chartData.length)
  isCreatingNewChart.value = true
  isChartLoaded.value = false
}

const vFocus = {
  mounted(el) {
    el.focus()
  }
}

defineExpose()

function quickPause(func) {
  setTimeout(func, 100)
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
        <DialogDescription>Load a chart from your saved charts</DialogDescription>
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
            <div class="w-8 invisible group-hover:visible flex justify-end" @click.stop="deleteChart(chart.id)">
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
        <div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm" @click="resetChart">
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
        <Input v-model="chartFull.title" placeholder="Enter chart title" />
      </div>

      <!-- type selection -->
      <div>
        <Label class="mb-1">Chart Type:</Label>
        <select v-model="chartFull.type" class="w-full p-2 border border-gray-300 rounded">
          <option v-for="type in chartTypes" :key="type">{{ type }}</option>
        </select>
      </div>

      <!-- axis naming -->
      <div>
        <Label class="mb-1">X Axis Name</Label>
        <Input placeholder="Enter X Axis name" v-model="chartFull.xColumn" />
        <Label class="mb-1 mt-4">Y Axis Name</Label>
        <Input placeholder="Enter Y Axis name" v-model="chartFull.yColumn" />
      </div>

      <!-- series management -->
      <div class="flex justify-between items-center">
        <p class="leading-7 [&:not(:first-child)]:mt-6">Series</p>
        <Plus @click="addSeriesSection" class="cursor-pointer"></Plus>
      </div>

      <!-- series scrollable container -->
      <div class="flex-1 overflow-y-scroll gap-4 flex flex-col">
        <div v-for="(section, index) in chartFull.chartData" :key="index"
          class="min-w-[300px] shrink-0 border border-gray-300 p-4 rounded-lg bg-gray-50 relative">
          <div class="flex justify-between items-center mb-2">
            <div class="flex items-center gap-2" v-if="!section.isEditingName ?? false"
              @click="section.isEditingName = true; section.editName = section.name;">
              <h2 class="text-md font-semibold text-gray-700">{{ section.name }}</h2>
              <Pencil size="16" class="text-neutral-600"></Pencil>
            </div>
            <div v-else class="flex items-center gap-2">
              <Input :id="`edit-series-name-${index}`" placeholder="Enter series name" v-model="section.editName"
                v-focus @keydown.enter="section.isEditingName = false; section.name = section.editName"
                @blur="quickPause(() => { section.isEditingName = false });" />
              <Check size="16" class="cursor-pointer"
                @click.stop="section.isEditingName = false; section.name = section.editName"></Check>
            </div>
            <!-- Remove Button -->
            <button @click="removeChartSection(index)" class="text-red-500 hover:text-red-700 cursor-pointer text-sm">
              ✕
            </button>
          </div>

          <!-- Query -->
          <div class="mb-2">
            <label class="block text-gray-600 text-sm mb-1">Saved Query:</label>
            <select v-model="section.query" class="w-full p-1 border border-gray-300 rounded text-sm"
              @change="generateChart(index)">
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
        </div>

      </div>

    </div>
    <Separator orientation="vertical"></Separator>
    <!-- Render Chart -->
    <div class="flex-1 ml-4 my-auto" v-if="chartComponent && chartFull.chartData.length">
      <div v-if="true">
        <component :is="chartComponent" :chart="chartFull" :xAxisTitle="chartFull.xColumn" :yAxisTitle="chartFull.yColumn" />
      </div>
    </div>
  </div>
</template>
