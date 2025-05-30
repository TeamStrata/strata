<script setup>
import { ref, computed } from 'vue';

// Accept options as a prop
const props = defineProps({
  options: {
    type: Array,
    required: true,
  }
});

const activeValues = ref([]);
const isDropdownOpen = ref(false);
const searchQuery = ref('');

// Use the passed-in prop instead of a local array
const filteredOptions = computed(() => {
  return props.options.filter(
    option =>
      option.toLowerCase().includes(searchQuery.value.toLowerCase()) &&
      !activeValues.value.includes(option)
  );
});

function toggleDropdown() {
  isDropdownOpen.value = !isDropdownOpen.value;
}

function selectOption(option) {
  activeValues.value.push(option);
  searchQuery.value = '';
}

function removeOption(option) {
  activeValues.value = activeValues.value.filter(val => val !== option);
}
</script>

<template>
  <div class="relative w-64">
    <div class="flex flex-wrap gap-2 mb-2">
      <span v-for="val in activeValues" :key="val"
        class="bg-blue-100 outline-1 text-blue-800 px-2 py-1 rounded-full flex items-center gap-1 text-sm">
        {{ val }}
        <button @click="removeOption(val)" class="text-blue-600 hover:text-red-600 cursor-pointer">×</button>
      </span>
    </div>

    <div class="relative">
      <button @click="toggleDropdown" class="px-4 py-2 rounded text-neutral-600 cursor-pointer w-full text-left">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 50 50">
          <path fill="currentColor"
            d="M25 42c-9.4 0-17-7.6-17-17S15.6 8 25 8s17 7.6 17 17s-7.6 17-17 17m0-32c-8.3 0-15 6.7-15 15s6.7 15 15 15s15-6.7 15-15s-6.7-15-15-15" />
          <path fill="currentColor" d="M16 24h18v2H16z" />
          <path fill="currentColor" d="M24 16h2v18h-2z" />
        </svg>
      </button>

      <div v-if="isDropdownOpen" class="absolute z-10 mt-2 w-full bg-white rounded shadow-lg">
        <input type="text" v-model="searchQuery" placeholder="Search..."
          class="w-full px-3 py-2 border-b outline-none" />
        <ul class="max-h-40 overflow-auto">
          <li v-for="option in filteredOptions" :key="option" @click="selectOption(option)"
            class="px-3 py-2 cursor-pointer hover:bg-blue-100">
            {{ option }}
          </li>
          <li v-if="filteredOptions.length === 0" class="px-3 py-2 text-gray-500">No options</li>
        </ul>
      </div>
    </div>
  </div>
</template>