<script setup>
import { apiFetch } from '@/api/request';
import router from '@/router';
import { useUserStore } from '@/stores/user';

const store = useUserStore()

function logout() {

    store.username = "";
    localStorage.removeItem('username');
    apiFetch("/logout", 'POST').then((res) => {
        if (res.ok) {
            router.push("/");
        } else {
            console.error("Something went wrong while logging out")
        }
    })

}

</script>


<template>
    <div class="flex flex-row justify-between bg-white border-1 border-gray-300 rounded-md p-2 px-4">
        <div>{{ store.username }}</div>

        <button @click="logout" class="cursor-pointer">Logout</button>
    </div>
</template>