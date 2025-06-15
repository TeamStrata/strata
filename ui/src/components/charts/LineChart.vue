<template>
  <apexchart
    type="line"
    height="350"
    :options="chartOptions"
    :series="chartSeries"
  />
</template>

<script setup>
import { computed } from 'vue';
import ApexCharts from 'vue3-apexcharts';

const props = defineProps({
  queryData: Array,
});

// Transform the queryData for ApexCharts
const chartSeries = computed(() => [
  {
    name: 'Value',
    data: props.queryData.map(item => item.value),
  },
]);

const chartOptions = computed(() => ({
  chart: {
    id: 'line-chart',
  },
  xaxis: {
    categories: props.queryData.map(item => item.name),
  },
}));
</script>