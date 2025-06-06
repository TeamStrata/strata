<template>
  <apexchart
    type="area"
    height="350"
    :options="chartOptions"
    :series="chartSeries"
  />
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  queryData: Array,
});

const chartSeries = computed(() => [
  { name: 'Value', data: props.queryData.map(item => item.value) },
]);

const chartOptions = computed(() => ({
  chart: { id: 'area-chart' },
  dataLabels: { enabled: false },
  stroke: { curve: 'smooth' },
  xaxis: { categories: props.queryData.map(item => item.name) },
}));
</script>