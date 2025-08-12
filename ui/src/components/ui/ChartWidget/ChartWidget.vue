<template>
  <div class="chart-widget" :style="style" ref="widgetRef">
    <div class="header">
      <h3>{{ chart.title }}</h3>
      <button v-if="props.editMode" @click="$emit('close')" class="cursor-pointer">
        <X></X>
      </button>

    </div>
    <div class="chart-container">
      <component :is="chartComponent" :chart="props.chart" :xAxisTitle="props.chart.xname"
        :yAxisTitle="props.chart.yname" />
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch, onBeforeUnmount, reactive, computed, defineAsyncComponent } from 'vue';
import { X } from 'lucide-vue-next';

const props = defineProps({
  chart: {
    type: Object,
    required: true
  },
  width: { type: Number, default: 300 },
  height: { type: Number, default: 300 },
  editMode: Boolean
});

//decide what chart to render
const chartComponent = computed(() => {
  const map = {
    line: defineAsyncComponent(() => import('../../charts/LineChart.vue')),
    area: defineAsyncComponent(() => import('../../charts/AreaChart.vue')),
    column: defineAsyncComponent(() => import('../../charts/ColumnChart.vue')),
    bar: defineAsyncComponent(() => import('../../charts/BarChart.vue')),
    scatter: defineAsyncComponent(() => import('../../charts/ScatterChart.vue')),
  }

  return map[props.chart.type] || null
})

const widgetWidth = ref(props.chart.size_x || 800);
const widgetHeight = ref(props.chart.size_y || 600);

const emit = defineEmits(['close', 'update:size', 'update']);

const widgetRef = ref(null);
let chartInstance = null;
let resizeObserver = null;

const style = computed(() => ({
  width: widgetWidth.value + 'px',
  height: widgetHeight.value + 'px',
  // resize: props.editMode ? 'both' : 'none',
  overflow: 'auto',
  border: '1px solid #ccc',
  borderRadius: '8px',
  padding: '1rem',
  backgroundColor: 'white',
  width: '100%',
  height: '100%',
  // minWidth: '300px',
  // minHeight: '300px',
  // maxWidth: '800px',
  // maxHeight: '600px'
}));

const observeSize = () => {
  const observer = new ResizeObserver(entries => {
    if (!props.editMode) return
    for (const entry of entries) {
      const { width, height } = entry.contentRect;
      widgetWidth.value = Math.round(width);
      widgetHeight.value = Math.round(height);

      // You can still emit if you want to notify parent
      emit('update:size', {
        id: props.chart.id,
        width: widgetWidth.value,
        height: widgetHeight.value,
      });
    }
  });
  if (widgetRef.value) observer.observe(widgetRef.value);
};


onMounted(() => {
  // fetchSeriesAndRender();
  observeSize();
});

watch(() => [widgetWidth.value, widgetHeight.value], () => {
  emit('update', {
    id: props.chart.id,
    width: widgetWidth.value,
    height: widgetHeight.value,
    x: props.chart.x,
    y: props.chart.y
  });
});
onBeforeUnmount(() => {
  if (chartInstance) {
    chartInstance.destroy();
  }
  if (resizeObserver && widgetRef.value) {
    resizeObserver.unobserve(widgetRef.value);
  }
});
</script>

<style scoped>
.chart-widget {
  position: relative;
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 8px rgba(0, 0, 0, 0.1);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.chart-container {
  flex-grow: 1;
  height: 400px;
  min-height: 300px;
}

canvas {
  width: 100% !important;
  height: 100% !important;
}
</style>
