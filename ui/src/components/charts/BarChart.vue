<template>
  <apexchart
    type="bar"
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
  chart: { id: 'bar-chart' },
  plotOptions: { bar: { horizontal: true } },
  xaxis: { categories: props.queryData.map(item => item.name) },
}));
</script>