<template>
  <div class="dashboard-container">
    <!-- Load Saved Charts Box -->
    <div class="bg-white border border-gray-300 rounded-xl p-4 shadow mb-6 w-full max-w-md">
      <h3 class="text-lg font-semibold text-gray-800 mb-2">Saved Chart</h3>
      <label class="block text-gray-600 font-medium mb-1">Select from saved charts:</label>
      <select v-model="selectedChartTitle" @change="loadChartFromDB" class="w-full p-2 border border-gray-300 rounded">
        <option value="">-- Select a chart --</option>
        <option v-for="chart in savedChartTitles" :key="chart.id" :value="chart">
          {{ chart.title }}
        </option>
      </select>
    </div>

    <!-- Selected chart(s) shown below -->
    <div class="chart-display">
      <ChartWidget v-for="chart in selectedCharts" :key="chart.id" :chart-data="chart" :width="chart.size_x"
        :height="chart.size_y" :id="`chart-widget-${chart.id}`" @close="removeChart(chart.id)"
        @update:size="(event) => onSizeUpdate(chart.id, event)" />
    </div>
    <br />
    <div>
      <button @click="saveDashboardCharts"
        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition">Save Layout</button>
    </div>
    <Toast ref="toastRef" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { apiFetch } from '@/api/request';
import { useRoute } from 'vue-router';
import Toast, { ToastTypes } from './Toast.vue';
import ChartWidget from './ui/ChartWidget/ChartWidget.vue';
import Swal from 'sweetalert2';
import { usePageInfoStore } from '@/stores/pageInfo';
const route = useRoute();
const dashboardId = ref(route.params.id);

const charts = ref([]);
const selectedCharts = ref([]);
const selectedChartId = ref("");
const selectedChartTitle = ref("");
const toastRef = ref(null);
const savedChartTitles = ref([])

const pageInfoStore = usePageInfoStore();

function onSizeUpdate(chartId, size) {
  const chart = selectedCharts.value.find(c => c.id === chartId);
  if (chart) {
    chart.size_x = size.width;
    chart.size_y = size.height;
  }
}

const fetchSavedChartTitles = async () => {
  try {
    const response = await apiFetch('/charts')
    if (!response.ok) throw new Error('Failed to fetch chart titles')
    savedChartTitles.value = await response.json()
  } catch (err) {
    console.error('Error fetching titles:', err)
  }
}

const saveDashboardCharts = async () => {
  try {
    // Fetch existing charts from backend
    const existingRes = await apiFetch(`/dashboard/${dashboardId.value}/charts`, 'GET');
    if (!existingRes.ok) throw new Error('Failed to fetch charts from dashboard');

    const existingCharts = await existingRes.json();

    // Map chart_id to size, cast to Number to ensure matching works
    const existingChartMap = new Map(
      Array.isArray(existingCharts)
        ? existingCharts
          .filter(chart => chart?.chart_id != null)
          .map(chart => [
            Number(chart.chart_id), // ✅ Ensure consistent type
            { size_x: chart.size_x, size_y: chart.size_y }
          ])
        : []
    );

    for (let index = 0; index < selectedCharts.value.length; index++) {
      const chart = selectedCharts.value[index];
      if (!chart?.id) continue;

      const chartIdNum = Number(chart.id); // ✅ Normalize ID for map lookup

      // 🟡 Get actual rendered size from the DOM
      const widgetEl = document.getElementById(`chart-widget-${chart.id}`);
      let size_x = chart.size_x || 300;
      let size_y = chart.size_y || 300;

      if (widgetEl) {
        const rect = widgetEl.getBoundingClientRect();
        size_x = Math.round(rect.width);
        size_y = Math.round(rect.height);
      }

      const existing = existingChartMap.get(chartIdNum);

      const shouldUpdate =
        !existing ||
        existing.size_x !== size_x ||
        existing.size_y !== size_y;

      // Debug logs
      console.log(`Chart ID ${chart.id}:`);
      console.log(`  Existing:`, existing);
      console.log(`  New Size: ${size_x}x${size_y}`);
      console.log(`  shouldUpdate:`, shouldUpdate);

      if (!shouldUpdate) {
        console.log(`Chart ID ${chart.id} already exists with same size. Skipping...`);
        continue;
      }

      // Payload for PATCH
      const payload = {
        chart_order: chart.chart_order || index + 1,
        id: chart.id,
        size_x,
        size_y,
        title: chart.title || "",
        type: chart.type || ""
      };

      // Save or update chart
      const response = await apiFetch(
        `/dashboard/${dashboardId.value}/chart/${chart.id}`,
        'PATCH',
        JSON.stringify(payload)
      );

      if (!response.ok) {
        toastRef.value?.showToast(
          `Failed to save chart ${chart.name || chart.id} to dashboard`,
          ToastTypes.FAIL
        );
        return;
      }
    }

    toastRef.value?.showToast('Dashboard layout saved successfully', ToastTypes.SUCCESS);
  } catch (error) {
    console.error(error);
    toastRef.value?.showToast('Error saving dashboard layout', ToastTypes.FAIL);
  }
};

const loadDashboardGraphs = async () => {
  try {
    const response = await apiFetch(`/dashboard/${dashboardId.value}/charts`);

    if (response.status === 404) {
      return;
    }

    const graphMappings = await response.json();
    console.log(graphMappings);
    if (graphMappings != null) {
      for (const mapping of graphMappings) {
        const chartResponse = await apiFetch(`/chart/${mapping.chart_id}`);
        if (chartResponse.ok) {
          const chart = await chartResponse.json();
          chart.size_x = mapping.size_x;
          chart.size_y = mapping.size_y;
          chart.order = mapping.order;
          selectedCharts.value.push(chart);
        }
      }
    } else {
      selectedCharts.value = [];
    }
  } catch (err) {
    if (err?.response?.status !== 404) {
      toastRef.value?.showToast('Error loading dashboard graphs', ToastTypes.FAIL);
    }
    console.error(err);
  }
};

const removeChart = async (chartId) => {
  const result = await Swal.fire({
    title: 'Remove Chart?',
    text: 'Are you sure you want to remove this chart from the dashboard?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#d33',
    cancelButtonColor: '#3085d6',
    confirmButtonText: 'Yes, remove it!',
  });

  if (!result.isConfirmed) return;

  try {
    // Remove from the local state
    selectedCharts.value = selectedCharts.value.filter(chart => chart.id !== chartId);

    // Delete from backend
    const res = await apiFetch(`/dashboard/${dashboardId.value}/chart/${chartId}`, 'DELETE');
    if (!res.ok) {
      toastRef.value?.showToast(`Failed to remove chart ${chartId} from dashboard`, ToastTypes.FAIL);
    } else {
      toastRef.value?.showToast(`Chart ${chartId} removed from dashboard`, ToastTypes.SUCCESS);
    }
  } catch (err) {
    console.error(err);
    toastRef.value?.showToast('Error removing chart from dashboard', ToastTypes.FAIL);
  }
};


const loadChartFromDB = () => {
  if (!selectedChartTitle.value) return;
  const existing = selectedCharts.value.find(c => c.id === selectedChartTitle.value.id);
  if (existing) return;
  selectedCharts.value.push({
    ...selectedChartTitle.value,
    size_x: selectedCharts.value.size_x || 600,
    size_y: selectedChartTitle.value.size_y || 800,
    chart_order: selectedCharts.value.length + 1
  });
  selectedChartTitle.value = "";
};

import { watch } from 'vue';

watch(() => route.params.id,
  (newVal, old) => {
    dashboardId.value = newVal;
    loadDashboardGraphs();
    fetchSavedChartTitles();
  });

async function loadDashboard() {
  const response = await apiFetch(`/dashboard/${dashboardId.value}`);
  if (!response.ok) return;
  const data = await response.json();
  console.log(data);

  // Update store with API info
  pageInfoStore.setPageInfo(
    data.title || "",
    data.content || "",
    data.configurable ?? true
  );
}

onMounted(async () => {
  await loadDashboardGraphs();
  await fetchSavedChartTitles();

});
loadDashboard();
</script>

<style scoped>
.dashboard-container {
  padding: 1rem;
  max-width: 1200px;
  margin: 0 auto;
}

.charts-list {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 2rem;
}

.chart-card,
.create-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 1rem;
  min-width: 200px;
  background: white;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  align-items: start;
  gap: 0.5rem;
}

.create-card {
  border: 2px dashed #ccc;
  cursor: pointer;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.chart-display {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.tag {
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.8rem;
}
</style>
