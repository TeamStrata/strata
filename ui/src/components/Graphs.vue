<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

// Graph types
const graphTypes = ['Line', 'Area', 'Column', 'Bar', 'Scatter'];
const selectedGraph = ref(graphTypes[0]);

// Saved queries
const savedQueries = ref([]);
const selectedQuery = ref(null);

// Fetch saved queries from backend
onMounted(async () => {
  try {
    const response = await axios.get('/query/list');
    savedQueries.value = response.data;
    selectedQuery.value = savedQueries.value[0]; // Preselect the first query
  } catch (error) {
    console.error('Error fetching saved queries:', error);
  }
});
</script>

<template>
  <div class="max-w-md mx-auto mt-10 p-6 bg-white shadow-md rounded-lg">
    <h1 class="text-2xl font-bold mb-4 text-gray-800">Graph Viewer</h1>

    <!-- Graph Type Dropdown -->
    <div class="flex items-center gap-4 mb-4">
      <label for="graphType" class="text-gray-600 font-medium">Graph Type:</label>
      <select
        id="graphType"
        v-model="selectedGraph"
        class="p-2 border border-gray-300 rounded focus:outline-none focus:ring focus:ring-blue-400"
      >
        <option v-for="type in graphTypes" :key="type">{{ type }}</option>
      </select>
    </div>

    <!-- Saved Query Dropdown -->
    <div class="flex items-center gap-4 mb-6">
      <label for="querySelect" class="text-gray-600 font-medium">Saved Query:</label>
      <select
        id="querySelect"
        v-model="selectedQuery"
        class="p-2 border border-gray-300 rounded focus:outline-none focus:ring focus:ring-green-400"
      >
        <option
          v-for="query in savedQueries"
          :key="query"
        >
          {{ query }}
        </option>
      </select>
    </div>

    <!-- Graph Display Area -->
    <div
      class="p-4 border border-dashed border-blue-300 rounded-lg flex flex-col items-center justify-center transition-colors duration-300"
    >
      <p class="text-lg font-semibold text-blue-700">Graph Type: {{ selectedGraph }}</p>
      <p class="text-lg font-semibold text-green-700">Query: {{ selectedQuery }}</p>

      <div
        v-if="selectedGraph && selectedQuery"
        class="mt-4 text-gray-800 font-medium"
      >
        <!-- Here, you can render the actual graph component! -->
        <p>Showing <span class="font-bold">{{ selectedGraph }}</span> graph for <span class="font-bold">{{ selectedQuery }}</span></p>
      </div>
    </div>
  </div>
</template>