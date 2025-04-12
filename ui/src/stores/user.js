import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  const username = ref("")

  function $reset() {
    username.value = ""
  }
  
  return { username }
})
