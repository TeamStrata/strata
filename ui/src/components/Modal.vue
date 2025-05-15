<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
    show: {
        type: Boolean,
        required: true,
    }
})

const emit = defineEmits(['close'])
</script>


<template>
  <!-- Backdrop Transition -->
  <Transition name="fade">
    <div
      v-if="show"
      class="fixed inset-0 bg-black/30 backdrop-blur-xs z-40"
      @click="emit('close')"
    />
  </Transition>

  <!-- Modal Content Transition -->
  <Transition name="scale">
    <div
      v-if="show"
      class="fixed z-50 left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2"
    >
      <div class="bg-white p-8 rounded-lg shadow-md">
        <slot />
      </div>
    </div>
  </Transition>
</template>
<style scoped>
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.scale-enter-active,
.scale-leave-active {
    transition: transform 0.15s ease, opacity 0.15s ease;
}

.scale-enter-from,
.scale-leave-to {
    transform: scale(0.95);
    opacity: 0;
}
</style>