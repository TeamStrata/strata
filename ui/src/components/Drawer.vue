<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
    show: {
        type: Boolean,
        required: true,
    }
});

const emit = defineEmits(['close']);
</script>

<template>
    <Transition name="fade">
        <div v-if="show" class="fixed inset-0 bg-black/30 backdrop-blur-sm z-40" @click="emit('close')"></div>
    </Transition>

    <Transition name="slide">
        <div v-if="show"
            class="fixed top-0 right-0 h-full w-80 bg-white shadow-xl z-50 p-5 overflow-auto">
            <!-- content goes here -->
            <slot />
        </div>
    </Transition>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
    transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
    transform: translateX(100%);
}
</style>
