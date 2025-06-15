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
  chart: { id: 'column-chart' },
  plotOptions: { bar: { columnWidth: '50%' } },
  xaxis: { categories: props.queryData.map(item => item.name) },
}));
</script>