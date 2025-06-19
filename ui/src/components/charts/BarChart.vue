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
})

const chartData = computed(() => {
  const firstSection = props.series[0]
  const labels = firstSection?.chartData.map(row => row[firstSection.xColumn]) ?? []

  const datasets = props.series.map((section, index) => ({
    label: section.query?.name || `Series ${index + 1}`,
    data: section.chartData.map(row => {
      const val = row?.[section.yColumn]
      return isNaN(val) ? 0 : Number(val)
    }),
    backgroundColor: `hsl(${(index * 60) % 360}, 70%, 60%)`
  }))

  return {
    labels,
    datasets,
  }
})

const chartOptions = {
  indexAxis: 'y',
  responsive: true,
  plugins: {
    legend: {
      display: true,
      position: 'top',
    },
    title: {
      display: true,
      text: 'Bar Chart',
    },
  },
}
</script>