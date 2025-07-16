import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', () => {

  //values that I need to store

  //actual code
  const code = ref("");
  //tab name (set to saved name?)
  const name = ref("");

  return { }
})