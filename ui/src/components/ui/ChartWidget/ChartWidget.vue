<template>
  <div class="chart-widget" :style="style" ref="widgetRef">
    <div class="header">
      <h3>{{ chartData.title }}</h3>
      <button v-if="props.editMode" @click="$emit('close')" class="cursor-pointer">
        <X></X>
      </button>

    </div>
    <div class="chart-container">
      <canvas ref="chartCanvas"></canvas>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch, onBeforeUnmount, reactive, computed } from 'vue';
import Chart from 'chart.js/auto';
import { apiFetch } from '@/api/request';
import { X } from 'lucide-vue-next';

const props = defineProps({
  chartData: {
    type: Object,
    required: true
  },
  width: { type: Number, default: 300 },
  height: { type: Number, default: 300 },
  editMode: Boolean
});


const widgetWidth = ref(props.width || 800);
const widgetHeight = ref(props.height || 600);

const emit = defineEmits(['close', 'update:size', 'update']);

const chartCanvas = ref(null);
const widgetRef = ref(null);
let chartInstance = null;
let resizeObserver = null;

const style = computed(() => ({
  width: widgetWidth.value + 'px',
  height: widgetHeight.value + 'px',
  resize: props.editMode ? 'both' : 'none',
  overflow: 'auto',
  border: '1px solid #ccc',
  borderRadius: '8px',
  padding: '1rem',
  backgroundColor: 'white',
  minWidth: '300px',
  minHeight: '300px',
  maxWidth: '800px',
  maxHeight: '600px'
}));

function getRandomColor() {
  const r = Math.floor(Math.random() * 200);
  const g = Math.floor(Math.random() * 200);
  const b = Math.floor(Math.random() * 200);
  return `rgb(${r}, ${g}, ${b})`;
}

const fetchSeriesAndRender = async () => {
  if (chartInstance) chartInstance.destroy();

  try {
    const chartRes = await apiFetch(`/chart/${props.chartData.id}`);
    if (!chartRes.ok) {
      console.error('Failed to load chart');
      return;
    }

    const chartInfo = await chartRes.json();
    const chartTypeRaw = chartInfo.type || 'line';
    const chartType = chartTypeRaw === 'column' ? 'bar'
      : chartTypeRaw === 'area' ? 'line'
        : chartTypeRaw;

    const response = await apiFetch(`/chart/${props.chartData.id}/series`);
    if (!response.ok) {
      console.error('Failed to load chart series');
      return;
    }

    const seriesList = await response.json();
    if (!seriesList.length) return;

    const datasets = await Promise.all(seriesList.map(async (series) => {
      const result = await apiFetch(`/query/${series.query_id}/execute`);
      if (!result.ok) return null;
      const data = await result.json();

      return {
        label: `${series.x_col_name}`,
        type: chartTypeRaw === 'area' ? 'line' : undefined,
        data: data.map(row => ({
          x: row[series.x_col_name],
          y: row[series.y_col_name]
        })),
        // backgroundColor: chartTypeRaw === 'area' ? getRandomColor(0.3) : getRandomColor(),
        // borderColor: getRandomColor(),
        fill: chartTypeRaw === 'area',
        tension: chartTypeRaw === 'area' ? 0.4 : 0,
        pointRadius: chartTypeRaw === 'scatter' ? 3 : 0,
        showLine: chartTypeRaw !== 'scatter'
      };
    }));

    const filteredDatasets = datasets.filter(Boolean);

    const allX = filteredDatasets.flatMap(d => d.data.map(p => Number(p.x))).filter(v => !isNaN(v));
    const allY = filteredDatasets.flatMap(d => d.data.map(p => Number(p.y))).filter(v => !isNaN(v));

    const xLabel = seriesList[0]?.x_col_name || 'X Axis';
    const yLabel = seriesList[0]?.y_col_name || 'Y Axis';

    chartInstance = new Chart(chartCanvas.value, {
      type: chartType,
      data: {
        datasets: filteredDatasets
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        parsing: false,
        scales: {
          x: {
            type: (chartTypeRaw === 'area' || chartTypeRaw === 'bar') && typeof filteredDatasets[0]?.data?.[0]?.x === 'string'
              ? 'category'
              : ['line', 'area', 'scatter', 'bar', 'column'].includes(chartTypeRaw)
                ? 'linear'
                : 'category',
            suggestedMin: allX.length ? Math.min(...allX) - 1 : 0,
            suggestedMax: allX.length ? Math.max(...allX) + 1 : 5,
            title: {
              display: true,
              text: xLabel
            },
            grid: {
              display: true
            }
          },
          y: {
            beginAtZero: true,
            suggestedMin: allY.length ? Math.min(...allY) - 1 : 0,
            suggestedMax: allY.length ? Math.max(...allY) + 1 : 5,
            title: {
              display: true,
              text: yLabel
            },
            ticks: {
              stepSize: 1,
              precision: 0
            },
            grid: {
              display: true
            }
          }
        },
        plugins: {
          legend: {
            position: 'top'
          },
          tooltip: {
            mode: 'nearest',
            intersect: false
          }
        }
      }
    });
  } catch (err) {
    console.error('Error rendering chart', err);
  }
};
const chartContainer = ref(null);

const observeSize = () => {
  const observer = new ResizeObserver(entries => {
    if(!props.editMode) return
    for (const entry of entries) {
      const { width, height } = entry.contentRect;
      widgetWidth.value = Math.round(width);
      widgetHeight.value = Math.round(height);

      // You can still emit if you want to notify parent
      emit('update:size', {
        id: props.chartData.id,
        width: widgetWidth.value,
        height: widgetHeight.value,
      });
    }
  });
  if (widgetRef.value) observer.observe(widgetRef.value);
};


onMounted(() => {
  fetchSeriesAndRender();
  observeSize();
});

watch(() => [widgetWidth.value, widgetHeight.value], () => {
  emit('update', {
    id: props.chartData.id,
    width: widgetWidth.value,
    height: widgetHeight.value,
    x: props.chartData.x,
    y: props.chartData.y
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
