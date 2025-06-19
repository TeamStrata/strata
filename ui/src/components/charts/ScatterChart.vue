<template>
  <apexchart
    type="line"
    height="350"
    :options="chartOptions"
    :series="chartSeries"
  />
</template>

<script setup>
import { computed } from 'vue'
import ApexCharts from 'vue3-apexcharts'

// Props: data array, x field name, y field name
const props = defineProps({
  queryData: {
    type: Array,
    required: true,
  },
  x: {
    type: String,
    required: true,
  },
  y: {
    type: String,
    required: true,
  },
})

const chartSeries = computed(() => [
  {
    name: props.y,
    data: props.queryData.map((item) => item[props.y]),
  },
])

const chartOptions = computed(() => ({
  chart: {
    id: 'line-chart',
  },
  xaxis: {
    categories: props.queryData.map((item) => item[props.x]),
  },
}))
</script>