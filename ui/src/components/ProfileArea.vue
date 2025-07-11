<script setup>
import { apiFetch } from '@/api/request';
import router from '@/router';
import { useUserStore } from '@/stores/user';

import DropdownMenuItem from './ui/dropdown-menu/DropdownMenuItem.vue';
import DropdownMenuLabel from './ui/dropdown-menu/DropdownMenuLabel.vue';
import DropdownMenuSeparator from './ui/dropdown-menu/DropdownMenuSeparator.vue';

const user = useUserStore()

function logout() {

    user.username = "";
    localStorage.removeItem('username');
    apiFetch("/logout", 'POST').then((res) => {
        if (res.ok) {
            router.push("/login");
        } else {
            console.error("Something went wrong while logging out")
        }
    })

}

</script>


<template>
    <!-- <div>{{ store.username }}</div> -->

    <!-- <button @click="logout" class="cursor-pointer">Logout</button> -->
    <DropdownMenuLabel>
        <div class="flex items-center gap-3">
            <span class="flex items-center justify-center bg-neutral-200 text-black rounded-full w-8 h-8 text-lg">
                {{ user.username[0].toUpperCase() }}
            </span>
            <span class="scroll-m-20 text-lg font-semibold tracking-tight">{{ user.username }}</span>
        </div>
    </DropdownMenuLabel>
    <DropdownMenuSeparator></DropdownMenuSeparator>
    <DropdownMenuItem :onclick="logout" class="flex justify-between hover:cursor-pointer">
        <span>Sign out</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="m16 17l5-5l-5-5m5 5H9m0 9H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
        </svg>
    </DropdownMenuItem>

</template>