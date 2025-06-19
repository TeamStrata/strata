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
  // Filler removed since no fill needed
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
  // Filler removed
)

const props = defineProps({
  series: {
    type: Array,
    required: true,
  }
})

// Build chart.js format data with multiple series
const chartData = computed(() => {
  return {
    labels: props.series[0]?.chartData?.map(item => item?.[props.series[0].xColumn]) ?? [],
    datasets: props.series.map((section, idx) => ({
      label: section.yColumn || `Series ${idx + 1}`,
      data: section.chartData.map(item => item?.[section.yColumn]),
      fill: false, // No fill, only line
      borderColor: `hsl(${idx * 60}, 70%, 50%)`,
      backgroundColor: `hsl(${idx * 60}, 70%, 50%)`, // backgroundColor can be same as border or removed
      tension: 0.4,
    })),
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
}
</script>