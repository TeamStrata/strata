import './assets/base.css'
import 'highlight.js/styles/atom-one-light.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// save session persistance
import { useUserStore } from './stores/user'
const user = useUserStore()

user.$subscribe((mut, state) => {
    console.log("CHANGE")
    console.log(mut, state)
    localStorage.setItem('user', JSON.stringify(state))
})

app.mount('#app')
