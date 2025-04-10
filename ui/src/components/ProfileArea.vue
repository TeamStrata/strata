<script setup>
import router from '@/router';


function logout() {
    router.push("/login");
}

</script>


<template>
    <div class="flex flex-row justify-between bg-white border-1 border-gray-300 rounded-md p-2 px-4">
        <div>Name</div>

        <button @click="logout" class="cursor-pointer">Logout</button>
    </div>
</template>