import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  const username = ref("")

  if (localStorage.getItem('user')) {
    const stored = JSON.parse(localStorage.getItem('user'))
    username.value = stored.username || "";
  }

  const isLoggedIn = computed(() => {
    return username.value != "";
  })
  return { username, isLoggedIn }
})
