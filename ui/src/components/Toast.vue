<script lang="ts" setup>
import { ref } from 'vue';


const show = ref(true)
const message = ref('Incorrect Toast Call')
const typeInternal = ref(ToastTypes.SUCCESS)

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
    <div v-if="show"
        class="fixed bottom-10 right-1/2 transform translate-x-1/2 px-4 py-2 rounded shadow border-l-6 flex items-center"
        :class="{
            'bg-green-50 border-green-400 text-green-500': typeInternal === ToastTypes.SUCCESS,
            'bg-amber-50 border-amber-400 text-amber-500': typeInternal === ToastTypes.WARNING,
            'bg-red-50 border-red-500 text-red-400': typeInternal === ToastTypes.FAIL,
            'bg-blue-50 border-blue-400 text-blue-400': typeInternal === ToastTypes.INFO,
        }">
        <svg v-if="typeInternal === ToastTypes.SUCCESS" class="mr-3 -ml-1" xmlns="http://www.w3.org/2000/svg" width="20"
            height="20" viewBox="0 0 24 24">
            <g fill="currentColor">
                <path d="M10.243 16.314L6 12.07l1.414-1.414l2.829 2.828l5.656-5.657l1.415 1.415z" />
                <path fill-rule="evenodd"
                    d="M1 12C1 5.925 5.925 1 12 1s11 4.925 11 11s-4.925 11-11 11S1 18.075 1 12m11 9a9 9 0 1 1 0-18a9 9 0 0 1 0 18"
                    clip-rule="evenodd" />
            </g>
        </svg>
        <svg v-else-if="typeInternal === ToastTypes.WARNING || ToastTypes.INFO" class="mr-3 -ml-1" xmlns="http://www.w3.org/2000/svg"
            width="20" height="20" viewBox="0 0 24 24">
            <g fill="currentColor">
                <path d="M11 10.98a1 1 0 1 1 2 0v6a1 1 0 1 1-2 0zm1-4.929a1 1 0 1 0 0 2a1 1 0 0 0 0-2" />
                <path fill-rule="evenodd"
                    d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10s10-4.477 10-10S17.523 2 12 2M4 12a8 8 0 1 0 16 0a8 8 0 0 0-16 0"
                    clip-rule="evenodd" />
            </g>
        </svg>
        <svg v-else-if="typeInternal === ToastTypes.FAIL" class="mr-3 -ml-1" xmlns="http://www.w3.org/2000/svg"
            width="20" height="20" viewBox="0 0 24 24">
            <g fill="none">
                <path
                    d="m12.593 23.258l-.011.002l-.071.035l-.02.004l-.014-.004l-.071-.035q-.016-.005-.024.005l-.004.01l-.017.428l.005.02l.01.013l.104.074l.015.004l.012-.004l.104-.074l.012-.016l.004-.017l-.017-.427q-.004-.016-.017-.018m.265-.113l-.013.002l-.185.093l-.01.01l-.003.011l.018.43l.005.012l.008.007l.201.093q.019.005.029-.008l.004-.014l-.034-.614q-.005-.018-.02-.022m-.715.002a.02.02 0 0 0-.027.006l-.006.014l-.034.614q.001.018.017.024l.015-.002l.201-.093l.01-.008l.004-.011l.017-.43l-.003-.012l-.01-.01z" />
                <path fill="currentColor"
                    d="M12 2c5.523 0 10 4.477 10 10s-4.477 10-10 10S2 17.523 2 12S6.477 2 12 2m0 2a8 8 0 1 0 0 16a8 8 0 0 0 0-16m0 11a1 1 0 1 1 0 2a1 1 0 0 1 0-2m0-9a1 1 0 0 1 1 1v6a1 1 0 1 1-2 0V7a1 1 0 0 1 1-1" />
            </g>
        </svg>
        {{ message }}

    </div>
</template>

<style></style>