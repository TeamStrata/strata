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
  LinearScale,
  TimeScale,
} from 'chart.js'
import { computed } from 'vue'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  TimeScale
)

const props = defineProps({
  series: {
    type: Array,
    required: true,
  },
})

const chartData = computed(() => ({
  datasets: props.series.map((s, index) => ({
    label: s.query?.name || `Series ${index + 1}`,
    data: (s.chartData || []).map((item) => ({
      x: item?.[s.xColumn],
      y: item?.[s.yColumn],
    })),
    borderColor: `hsl(${index * 60}, 70%, 50%)`,
    backgroundColor: `hsl(${index * 60}, 70%, 70%)`,
    tension: 0.4,
  })),
}))

const chartOptions = {
  responsive: true,
  parsing: false, // IMPORTANT: allows x/y object format
  scales: {
    x: {
      type: 'linear', // use 'time' if x is a date
      title: {
        display: true,
        text: 'X Axis',
      },
    },
    y: {
      title: {
        display: true,
        text: 'Y Axis',
      },
    },
  },
  plugins: {
    legend: {
      display: true,
    },
  },
}
</script>