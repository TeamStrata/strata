<template>
  <div class="chart-widget" :style="style">
    <div class="header">
      <h3>{{ chartData.name }}</h3>
      <button @click="$emit('close')">✖</button>
    </div>
    <div class="chart-container">
      <component
        v-if="ChartComponent && chartSeries.length"
        :is="ChartComponent"
        :series="chartSeries"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, shallowRef } from 'vue';
import { apiFetch } from '@/api/request';

const props = defineProps({
  chartData: {
    type: Object,
    required: true,
  },
});

const chartSeries = ref([]);
const ChartComponent = shallowRef(null);

const style = {
  resize: 'both',
  overflow: 'auto',
  border: '1px solid #ccc',
  borderRadius: '8px',
  padding: '1rem',
  backgroundColor: 'white',
  minWidth: '300px',
  minHeight: '300px',
  maxWidth: '600px',
  maxHeight: '600px',
};

const loadChartComponent = async (type) => {
  try {
    switch (type) {
      case 'line':
        ChartComponent.value = (await import('../../charts/LineChart.vue')).default;
        break;
      case 'bar':
        ChartComponent.value = (await import('../../charts/BarChart.vue')).default;
        break;
      case 'column':
        ChartComponent.value = (await import('../../charts/ColumnChart.vue')).default;
        break;
      case 'scatter':
        ChartComponent.value = (await import('../../charts/ScatterChart.vue')).default;
        break;
      case 'area':
        ChartComponent.value = (await import('../../charts/AreaChart.vue')).default;
        break;
      default:
        console.warn('Unknown chart type:', type);
        ChartComponent.value = null;
    }
  } catch (err) {
    console.error('Error loading chart component:', err);
  }
};

const fetchSeriesAndRender = async () => {
  try {
    const response = await apiFetch(`/chart/${props.chartData.id}/series`);
    if (!response.ok) {
      console.error('Failed to load chart series');
      return;
    }
    const seriesList = await response.json();

    const datasets = await Promise.all(
      seriesList.map(async (series) => {
        const result = await apiFetch(`/query/${series.query_id}/execute`);
        if (!result.ok) return null;
        const data = await result.json();

        return {
          ...series,
          chartData: data,
        };
      })
    );

    chartSeries.value = datasets.filter(Boolean);
  } catch (err) {
    console.error('Error fetching series:', err);
  }
};

watch(
  () => props.chartData,
  async (newChart) => {
    await loadChartComponent(newChart.chart_type);
    await fetchSeriesAndRender();
  },
  { immediate: true, deep: true }
);
</script>

<style scoped>
.chart-widget {
  position: relative;
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 8px rgba(0, 0, 0, 0.1);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.chart-container {
  flex-grow: 1;
  height: 100%;
}

canvas {
  width: 100% !important;
  height: 100% !important;
}

button {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  color: #888;
}
</style>
