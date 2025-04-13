<script lang="ts" setup>
import { ref } from 'vue';


const show = ref(false)
const message = ref('Incorrect Toast Call')
const typeInternal = ref(ToastTypes.WARNING)

function showToast(msg = 'Incorrect Toast Call', type = ToastTypes.WARNING, duration = 3000) {
    message.value = msg;
    typeInternal.value = type;

    show.value = true;

    setTimeout(() => {
        show.value = false;
    }, duration);
}

defineExpose({ showToast });
</script>

<script lang="ts">
export enum ToastTypes {
    SUCCESS,
    FAIL,
    WARNING,
    INFO
}
</script>

<template>
    <div v-if="show" class="fixed bottom-5 right-5 bg-black text-white px-4 py-2 rounded shadow">
        {{ message }}
    </div>
</template>
