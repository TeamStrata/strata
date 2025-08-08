import { markRaw } from 'vue'
import { defineStore } from 'pinia'

export const useEventBus = defineStore('eventBus', () => {
  const listeners = markRaw(new Map())

  function on(event, handler) {
    if (!listeners.has(event)) {
      listeners.set(event, new Set())
    }
    listeners.get(event).add(handler)
    return () => off(event, handler) // unsubscribe helper
  }

  function off(event, handler) {
    if (listeners.has(event)) {
      listeners.get(event).delete(handler)
    }
  }

  function emit(event, payload) {
    if (listeners.has(event)) {
      for (const fn of listeners.get(event)) {
        fn(payload)
      }
    }
  }

  return { on, off, emit }
})
