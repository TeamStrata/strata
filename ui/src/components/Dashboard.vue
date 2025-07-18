<template>
  <div class="dashboard-container">
    <h2>Charts</h2>

    <!-- List of saved charts -->
    <div class="charts-list">
      <div
        v-for="chart in charts"
        :key="chart.id"
        class="chart-card"
      >
        <strong>{{ chart.name }}</strong>
        <span class="tag" :style="{ backgroundColor: chart.color }">{{ chart.color }}</span>
        <button @click="editChart(chart)">Edit</button>
        <button @click="deleteChart(chart.id)">Delete</button>
      </div>

      <!-- Create a new chart -->
      <div class="create-card" @click="createChart">
        <span>＋ Create a New Chart</span>
      </div>
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
      <button @click="saveDashboardCharts">💾 Save Layout</button>
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
const toastRef = ref(null);

const loadCharts = async () => {
  try {
    const response = await apiFetch('/charts');
    if (!response.ok) {
      toastRef.value?.showToast('Failed to load charts', ToastTypes.FAIL);
      return;
    }
    charts.value = await response.json();
  } catch (err) {
    toastRef.value?.showToast('Error loading charts', ToastTypes.FAIL);
    console.error(err);
  }
};
const saveDashboardCharts = async () => {
  const dashboardGraphs = selectedCharts.value.map((chart, index) => ({
    dash_id: parseInt(dashboardId),
    chart_id: chart.id,
    size_x: chart.size_x || 300,
    size_y: chart.size_y || 300,
    order: chart.chart_order || index + 1
  }));

  try {
    const response = await apiFetch(`/dashboard`, {
      method: 'POST',
      body: JSON.stringify(dashboardGraphs),
      headers: {
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      toastRef.value?.showToast('Failed to save dashboard layout', ToastTypes.FAIL);
    } else {
      toastRef.value?.showToast('Dashboard layout saved successfully', ToastTypes.SUCCESS);
    }
  } catch (error) {
    console.error(error);
    toastRef.value?.showToast('Error saving dashboard layout', ToastTypes.FAIL);
  }
};
const loadDashboardGraphs = async () => {
  try {
    const response = await apiFetch(`/dashboard/${dashboardId}/charts`);
    if (!response.ok) {
      toastRef.value?.showToast('Failed to load dashboard graphs', ToastTypes.FAIL);
      return;
    }
    const graphMappings = await response.json();
    for (const mapping of graphMappings) {
      const chartResponse = await apiFetch(`/charts/${mapping.chart_id}`);
      if (chartResponse.ok) {
        const chart = await chartResponse.json();
        chart.size_x = mapping.size_x;
        chart.size_y = mapping.size_y;
        chart.chart_order = mapping.chart_order;
        selectedCharts.value.push(chart);
      }
    }
  } catch (err) {
    toastRef.value?.showToast('Error loading dashboard graphs', ToastTypes.FAIL);
    console.error(err);
  }
};

const editChart = (chart) => {
  if (!selectedCharts.value.find((c) => c.id === chart.id)) {
    selectedCharts.value.push(chart);
  }
};

const createChart = () => {
  const newChart = {
    id: Date.now(),
    name: 'New Chart',
    color: '#888',
    data: {},
    size_x: 300,
    size_y: 300,
    chart_order: selectedCharts.value.length + 1
  };
  selectedCharts.value.push(newChart);
};

const deleteChart = async (id) => {
  try {
    const response = await apiFetch(`/charts/${id}`, { method: 'DELETE' });
    if (!response.ok) {
      toastRef.value?.showToast('Delete failed', ToastTypes.FAIL);
      return;
    }
    charts.value = charts.value.filter((c) => c.id !== id);
    selectedCharts.value = selectedCharts.value.filter((c) => c.id !== id);
  } catch (err) {
    toastRef.value?.showToast('Error deleting chart', ToastTypes.FAIL);
    console.error(err);
  }
};

const removeChart = (id) => {
  selectedCharts.value = selectedCharts.value.filter((c) => c.id !== id);
};

onMounted(async () => {
  await loadCharts();
  await loadDashboardGraphs();
});
</script>

<style scoped>
.dashboard-container {
  padding: 1rem;
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
