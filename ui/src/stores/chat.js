import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat', () => {

    const messages = ref([])
    const initialized = ref(false)

  return { messages, initialized }
})
