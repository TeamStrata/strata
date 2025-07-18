<template>
  <div class="dashboard-container">
    <!-- Load Saved Charts Box -->
    <div class="bg-white border border-gray-300 rounded-xl p-4 shadow mb-6 w-full max-w-md">
      <h3 class="text-lg font-semibold text-gray-800 mb-2">Saved Chart</h3>
      <label class="block text-gray-600 font-medium mb-1">Select from saved charts:</label>
      <select
        v-model="selectedChartTitle"
        @change="loadChartFromDB"
        class="w-full p-2 border border-gray-300 rounded"
      >
        <option value="">-- Select a chart --</option>
        <option v-for="chart in savedChartTitles" :key="chart.id" :value="chart">
          {{ chart.title }}
        </option>
      </select>
    </div>

    <!-- Selected chart(s) shown below -->
    <div class="chart-display">
      <ChartWidget
        v-for="chart in selectedCharts"
        :key="chart.id"
        :chart-data="chart"
        @close="removeChart(chart.id)"
      />
    </div>

    <div>
      <button @click="saveDashboardCharts">Save Layout</button>
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

const route = useRoute();
const dashboardId = route.params.id;

const charts = ref([]);
const selectedCharts = ref([]);
const selectedChartId = ref("");
const selectedChartTitle = ref("");
const toastRef = ref(null);
const savedChartTitles = ref([])

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
    for (let index = 0; index < selectedCharts.value.length; index++) {
      const chart = selectedCharts.value[index];
      const payload = {
        size_x: chart.size_x || 300,
        size_y: chart.size_y || 300,
        order: chart.chart_order || index + 1
      };

      const response = await apiFetch(`/dashboard/${dashboardId}/chart/${chart.id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload),
        headers: {
          'Content-Type': 'application/json'
        }
      });

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
    const response = await apiFetch(`/dashboard/${dashboardId}/charts`);

    if (response.status === 404) {
      return; 
    }

    const graphMappings = await response.json();
    console.log(graphMappings);
    if(graphMappings!=null){
      for (const mapping of graphMappings) {
        const chartResponse = await apiFetch(`/charts/${mapping.chart_id}`);
        if (chartResponse.ok) {
          const chart = await chartResponse.json();
          chart.size_x = mapping.size_x;
          chart.size_y = mapping.size_y;
          chart.order = mapping.order;
          selectedCharts.value.push(chart);
        }
      }
    }
  } catch (err) {
    if (err?.response?.status !== 404) {
      toastRef.value?.showToast('Error loading dashboard graphs', ToastTypes.FAIL);
    }
    console.error(err);
  }
};

const removeChart = (id) => {
  selectedCharts.value = selectedCharts.value.filter((c) => c.id !== id);
};

const loadChartFromDB = () => {
  if (!selectedChartTitle.value) return;
  const existing = selectedCharts.value.find(c => c.id === selectedChartTitle.value.id);
  if (existing) return;
  selectedCharts.value.push({
    ...selectedChartTitle.value,
    size_x: 300,
    size_y: 300,
    chart_order: selectedCharts.value.length + 1
  });
  selectedChartTitle.value = "";
};

onMounted(async () => {
  await loadDashboardGraphs();
  await fetchSavedChartTitles();
});
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
