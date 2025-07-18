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

const renderChart = () => {
  if (chartInstance) {
    chartInstance.destroy();
  }

  const data = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
    datasets: [
      {
        label: props.chartData.name,
        data: [12, 19, 3, 5, 2, 3], // Placeholder, replace with real data
        backgroundColor: props.chartData.color || '#888'
      }
    ]
  };

  chartInstance = new Chart(chartCanvas.value, {
    type: 'bar',
    data,
    options: {
      responsive: true,
      maintainAspectRatio: false
    }
  });
};

onMounted(() => {
  renderChart();
});

watch(() => props.chartData, () => {
  renderChart();
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
