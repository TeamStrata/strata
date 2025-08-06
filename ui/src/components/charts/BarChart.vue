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

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

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

const chartData = computed(() => {
  const firstSection = props.series[0]
  const labels = firstSection?.chartData.map(row => row[firstSection.xColumn]) ?? []

  const datasets = props.series.map((section, index) => ({
    label: section.name || `Series ${index + 1}`,
      data: (section.chartData || []).map(item => ({
        y: item?.[section.xColumn],
        x: item?.[section.yColumn],
      })),
    backgroundColor: `hsl(${(index * 60) % 360}, 70%, 60%)`
  }))

  return {
    labels,
    datasets,
  }
})

const chartOptions = computed(() => {
  return {
    indexAxis: 'y',
    responsive: true,
    plugins: {
      legend: {
        display: true,
        position: 'top',
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