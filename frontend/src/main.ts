import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import { useAuthStore } from './stores/auth'

// Configure axios defaults
axios.defaults.baseURL = 'http://localhost:8080'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Initialize auth state before mounting
const authStore = useAuthStore()
await authStore.initializeAuth()

app.mount('#app')
