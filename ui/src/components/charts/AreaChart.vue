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
  Filler,
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
  Filler // Needed for area fill
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
      fill: true,
      borderColor: `hsl(${idx * 60}, 70%, 50%)`,
      backgroundColor: `hsla(${idx * 60}, 70%, 70%, 0.3)`,
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