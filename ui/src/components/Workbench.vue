<template>
  <div class="ai-workbench">
    <h1>AI Bot Workbench</h1>

    <textarea
      v-model="userInput"
      placeholder="Enter your prompt here..."
      rows="5"
      class="input-box"
    ></textarea>

    <button @click="sendQuery" :disabled="loading">
      {{ loading ? 'Sending...' : 'Send to AI' }}
    </button>

    <div v-if="response" class="response-box">
      <h3>AI Response:</h3>
      <pre>{{ response }}</pre>
    </div>

    <div v-if="error" class="error-box">
      <strong>Error:</strong> {{ error }}
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  setup() {
    const userInput = ref('')
    const response = ref('')
    const error = ref('')
    const loading = ref(false)

    const API_URL = '/api/ai-bot'
    const sendQuery = async () => {
      if (!userInput.value.trim()) return

      loading.value = true
      response.value = ''
      error.value = ''

      try {
        const res = await fetch(API_URL, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ prompt: userInput.value }),
        })

        if (!res.ok) throw new Error(`Error ${res.status}: ${res.statusText}`)

        const data = await res.json()
        response.value = data.response || 'No response from AI.'
      } catch (err) {
        error.value = err.message
      } finally {
        loading.value = false
      }
    }

    return {
      userInput,
      response,
      error,
      loading,
      sendQuery
    }
  }
}
</script>