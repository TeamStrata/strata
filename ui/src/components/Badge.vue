<script setup>
defineProps({
  color: {
    type: String,
    default: '#87d395', // fallback color
  },
})

</script>

<template>
  <div class="inline-block px-1 rounded-sm border text-sm" :style="{
    backgroundColor: color + '10',
    borderColor: color,
    color: color
  }">
    <slot></slot>
  </div>
</template>