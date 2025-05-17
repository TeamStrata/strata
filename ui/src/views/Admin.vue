<script setup>
import router from '@/router';
import { RouterLink, RouterView } from 'vue-router';


</script>

<template>
    <button class="absolute left-40 top-10 cursor-pointer" @click="router.push('/')">Back</button>
    <div class="flex flex-row max-w-[1100px] mx-auto h-screen p-20">
        <!-- sidebar -->
        <div class="w-3xs h-full p-5 flex flex-col">
            <RouterLink to="/admin/users">Users Button</RouterLink>
            <RouterLink to="/admin/roles">Roles Button</RouterLink>
        </div>
        <div class="h-full w-0.5 bg-neutral-200"></div>
        <div class="w-full h-full p-5">
            <!-- focus area -->
            <RouterView></RouterView>
        </div>
    </div>
</template>