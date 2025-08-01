<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { apiFetch } from './api/request';
import router from './router';

apiFetch("/auth", "POST").then(res => {
  if (res.status != 200) {
    router.push('/login');
  }
}).catch(error => {
  console.error("Error during authentication:", error);
})
</script>

<template>
  <main>
    <RouterView> </RouterView>
  </main>
</template>
