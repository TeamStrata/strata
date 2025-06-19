<template>
  <Bar :data="chartData" :options="chartOptions" />
</template>

<script setup>
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js'
import { computed } from 'vue'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
)

const props = defineProps({
  series: {
    type: Array,
    required: true,
  }
})

// Prepare chart data with multiple bar series
const chartData = computed(() => {
  const labels = props.series[0]?.chartData?.map(item => item?.[props.series[0].xColumn]) ?? []

  return {
    labels,
    datasets: props.series.map((section, idx) => ({
      label: section.yColumn || `Series ${idx + 1}`,
      data: section.chartData.map(item => item?.[section.yColumn]),
      backgroundColor: `hsl(${idx * 60}, 70%, 60%)`,
    }))
  }
})

const chartOptions = {
  responsive: true,
  plugins: {
    legend: {
      display: true,
    },
    title: {
      display: false,
    },
  },
  scales: {
    x: {
      title: {
        display: true,
        text: props.series[0]?.xColumn || 'X-Axis',
      }
    },
    y: {
      title: {
        display: true,
        text: 'Values',
      },
      beginAtZero: true
    }
  }
}
</script>