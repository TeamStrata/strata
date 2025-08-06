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
  },
  xAxisTitle: {
    type: String,
    required: true,
  },
  yAxisTitle: {
    type: String,
    required: true,
  },
})

// Prepare chart data with multiple bar series
const chartData = computed(() => {
  const labels = props.series[0]?.chartData?.map(item => item?.[props.series[0].xColumn]) ?? []

  return {
    labels,
    datasets: props.series.map((section, idx) => ({
      label: section.name || `Series ${idx + 1}`,
      data: section.chartData.map(item => ({
        x: item?.[section.xColumn],
        y: item?.[section.yColumn],
      })), backgroundColor: `hsl(${idx * 60}, 70%, 60%)`,
    }))
  }
})

const chartOptions = computed(() => {
  return {
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
        type: 'linear',
        title: {
          display: true,
          text: props.xAxisTitle,
        }
      },
      y: {
        title: {
          display: true,
          text: props.yAxisTitle,
        },
        beginAtZero: true
      }
    }
  }
})
</script>