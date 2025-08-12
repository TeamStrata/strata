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
  chart: {
    type: Object,
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

function labels() {
  const xCol = props.chart.xColumn;
  return props.chart.chartData[0].data.map(item => item[xCol]);
}


const chartData = computed(() => {
  return {
    // labels: props.series[0]?.chartData?.map(item => item?.[props.series[0].xColumn]) ?? [],
    labels: labels(),
    datasets: props.chart.chartData.map((section, idx) => ({
      label: section.name || `Series ${idx + 1}`,
      data: section.data.map(item => ({
        x: item?.[section.xColumn],
        y: item?.[section.yColumn],
      })),
      backgroundColor: `hsla(${idx * 60}, 70%, 70%)`,
      tension: 0.4
    })),
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