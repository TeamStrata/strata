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

// Build chart.js format data with multiple series
const chartData = computed(() => {
  return {
    labels: props.series[0]?.chartData?.map(item => item?.[props.series[0].xColumn]) ?? [],
    datasets: props.series.map((section, idx) => ({
      label: section.name || `Series ${idx + 1}`,
      data: section.chartData.map(item => ({
        x: item?.[section.xColumn],
        y: item?.[section.yColumn],
      })),      fill: true,
      borderColor: `hsl(${idx * 60}, 70%, 50%)`,
      backgroundColor: `hsla(${idx * 60}, 70%, 70%, 0.3)`,
      tension: 0.4
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
        text: props.xAxisTitle,
      },
    },
    y: {
      beginAtZero: false,
      title: {
        display: true,
        text: props.yAxisTitle,
      },
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