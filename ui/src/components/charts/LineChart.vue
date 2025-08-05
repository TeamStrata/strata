<template>
  <Line :data="chartData" :options="chartOptions" />
</template>

<script setup>
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
} from 'chart.js'
import { computed } from 'vue'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
)

const props = defineProps({
  series: {
    type: Array,
    required: true,
  },
})

const chartData = computed(() => {
  return {
    datasets: props.series.map((section, idx) => ({
      label: section.name || `Series ${idx + 1}`,
      data: section.chartData.map(item => ({
        x: item?.[section.xColumn],
        y: item?.[section.yColumn],
      })),
      fill: false,
      borderColor: `hsl(${idx * 60}, 70%, 50%)`,
      backgroundColor: `hsl(${idx * 60}, 70%, 50%)`,
      tension: 0.4,
    })),
  }
})
const chartOptions = {
  responsive: true,
  scales: {
    x: {
      type: 'linear',
      title: {
        display: true,
        text: props.series[0]?.xColumn || 'X Axis',
      },
      min: 0, 
    },
    y: {
      beginAtZero: true, 
    },
  },
  plugins: {
    legend: {
      display: true,
    },
    title: {
      display: false,
    },
  },
}
</script>