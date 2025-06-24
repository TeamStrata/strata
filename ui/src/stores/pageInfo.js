import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', () => {

  const title = ref("")
  const description = ref("")
  const configurable = ref(false)

  function setPageInfo(newTitle, newDesc, newConf) {
    title.value = newTitle
    description.value = newDesc
    configurable.value = newConf
  }

  const pageInfo = computed(() => {
    return {title, description, configurable }
  })

  return { pageInfo, setPageInfo }
})
