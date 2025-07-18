<template>
  <div class="chart-widget" :style="style">
    <div class="header">
      <h3>{{ chartData.name }}</h3>
      <button @click="$emit('close')">✖</button>
    </div>
    <div class="chart-container">
      <canvas ref="chartCanvas"></canvas>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue';
import Chart from 'chart.js/auto';
import { apiFetch } from '@/api/request';

const props = defineProps({
  chartData: {
    type: Object,
    required: true
  }
});

const chartCanvas = ref(null);
let chartInstance = null;

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
  maxHeight: '600px'
};

const fetchSeriesAndRender = async () => {
  if (chartInstance) {
    chartInstance.destroy();
  }

  try {
    const response = await apiFetch(`/chart/${props.chartData.id}/series`);
    if (!response.ok) {
      console.error('Failed to load chart series');
      return;
    }
    const seriesList = await response.json();

    const datasets = await Promise.all(seriesList.map(async (series) => {
      const result = await apiFetch(`/query/${series.query_id}/execute`);
      if (!result.ok) return null;
      const chartData = await result.json();

      return {
        label: series.name,
        data: chartData.map(row => ({
          x: row[series.x_col_name],
          y: row[series.y_col_name]
        })),
        fill: false
      };
    }));

    const filteredDatasets = datasets.filter(Boolean);

    chartInstance = new Chart(chartCanvas.value, {
      type: props.chartData.chart_type || 'line',
      data: {
        datasets: filteredDatasets
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        parsing: false,
        scales: {
          x: {
            type: 'category',
            title: {
              display: true,
              text: 'X Axis'
            }
          },
          y: {
            title: {
              display: true,
              text: 'Y Axis'
            },
            ticks: {
              callback: function(value) {
                return Number.isInteger(value) ? value : null;
              },
              stepSize: 1
            }
          }
        }
      }
    });
  } catch (err) {
    console.error('Error rendering chart', err);
  }
};

onMounted(() => {
  fetchSeriesAndRender();
});

watch(() => props.chartData, () => {
  fetchSeriesAndRender();
}, { deep: true });
</script>

<style scoped>
.chart-widget {
  position: relative;
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 8px rgba(0,0,0,0.1);
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
