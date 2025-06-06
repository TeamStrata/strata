<script setup>
import { ref, onMounted, watch } from 'vue';
import axios from 'axios';

// Import the graph components
import LineGraph from './graphs/LineGraph.vue';
import AreaGraph from './graphs/AreaGraph.vue';
import ColumnGraph from './graphs/ColumnGraph.vue';
import BarGraph from './graphs/BarGraph.vue';
import ScatterGraph from './graphs/ScatterGraph.vue';

const graphTypes = ['Line', 'Area', 'Column', 'Bar', 'Scatter'];
const selectedGraph = ref('Line');

const savedQueries = ref([]);
const selectedQuery = ref(null);
const queryData = ref([]);
const loading = ref(false);
const error = ref(null);

// Load saved queries from backend
onMounted(async () => {
  try {
    const response = await axios.get('/api/saved-queries');
    savedQueries.value = response.data;
    if (savedQueries.value.length > 0) {
      selectedQuery.value = savedQueries.value[0];
    }
  } catch (err) {
    console.error('Error fetching saved queries:', err);
    error.value = 'Failed to load saved queries';
  }
});

// Load query data whenever selectedQuery changes
watch(selectedQuery, async (newQuery) => {
  if (!newQuery) return;
  loading.value = true;
  error.value = null;
  try {
    const response = await axios.get(`/query/list?queryName=${encodeURIComponent(newQuery)}`);
    queryData.value = response.data;
  } catch (err) {
    console.error('Error fetching query data:', err);
    error.value = 'Failed to load query data';
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="max-w-2xl mx-auto mt-10 p-6 bg-white shadow-md rounded-lg">
    <h1 class="text-2xl font-bold mb-4 text-gray-800">Graph Viewer</h1>

    <!-- Graph Type Dropdown -->
    <div class="mb-4">
      <label class="text-gray-600 font-medium">Graph Type:</label>
      <select v-model="selectedGraph" class="ml-2 p-2 border rounded focus:outline-none focus:ring focus:ring-blue-400">
        <option v-for="type in graphTypes" :key="type">{{ type }}</option>
      </select>
    </div>

    <!-- Saved Query Dropdown -->
    <div class="mb-4">
      <label class="text-gray-600 font-medium">Query:</label>
      <select v-model="selectedQuery" class="ml-2 p-2 border rounded focus:outline-none focus:ring focus:ring-green-400">
        <option v-for="query in savedQueries" :key="query">{{ query }}</option>
      </select>
    </div>

    <!-- Loading & Error States -->
    <div v-if="loading" class="text-blue-600 font-medium mb-4">Loading data...</div>
    <div v-if="error" class="text-red-500 font-medium mb-4">{{ error }}</div>

    <!-- Graph Component -->
    <div v-if="queryData.length > 0" class="mt-6">
      <component
        :is="selectedGraph + 'Graph'"
        :query-data="queryData"
      />
    </div>

    <!-- No data fallback -->
    <div v-else-if="!loading && !error" class="text-gray-500 italic mt-4">
      No data available for this query.
    </div>
  </div>
</template>

<script>
// Register the graph components
export default {
  components: {
    LineGraph,
    AreaGraph,
    ColumnGraph,
    BarGraph,
    ScatterGraph,
  },
};
</script>