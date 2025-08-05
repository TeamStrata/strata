<template>
  <Scatter :data="chartData" :options="chartOptions" />
</template>

<script setup>
import { Scatter } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  PointElement,
  LinearScale,
  TimeScale,
} from 'chart.js'
import 'chartjs-adapter-date-fns' // date adapter for 'time' scale if needed
import { computed } from 'vue'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  PointElement,
  LinearScale,
  TimeScale,
)

const props = defineProps({
  series: {
    type: Array,
    required: true,
    // Each series: { chartData: Array, xColumn: String, yColumn: String }
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
  // Map each series in props to a Chart.js dataset
  return {
    datasets: props.series.map((section, idx) => ({
      label: section.name || `Series ${idx + 1}`,
      data: (section.chartData || []).map(item => ({
        x: item?.[section.xColumn],
        y: item?.[section.yColumn],
      })),
      backgroundColor: `hsl(${(idx * 60) % 360}, 70%, 50%)`,
      pointRadius: 5,
    })),
  }
})

const xIsDate = computed(() => {
  const firstSeries = props.series[0]
  if (!firstSeries || !firstSeries.chartData?.length) return false
  const sampleX = firstSeries.chartData[0]?.[firstSeries.xColumn]
  return typeof sampleX === 'string' && /\d{4}-\d{2}-\d{2}/.test(sampleX)
})

const chartOptions = computed(() => ({
  responsive: true,
  scales: {
    x: {
      type: xIsDate.value ? 'time' : 'linear',
      title: {
        display: true,
        text: props.xAxisTitle,
      },
      time: xIsDate.value
        ? {
          tooltipFormat: 'PP',
          unit: 'day',
        }
        : undefined,
    },
    y: {
      type: 'linear',
      title: {
        display: true,
        text: props.yAxisTitle,
      },
    },
  },
  plugins: {
    legend: {
      display: true,
      position: 'top',
    },
    tooltip: {
      enabled: true,
      mode: 'nearest',
      intersect: true,
    },
  },
}))
</script>
