<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { computed } from 'vue'

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
</script>

<template>
  <div class="navbar is-light" role="navigation" aria-label="main navigation">
    <nav class="navbar-menu">
      <div class="navbar-start">
        <router-link to="/" class="navbar-item">Home</router-link>
        <router-link to="/articles" class="navbar-item">Articles</router-link>
        <router-link to="/authors" class="navbar-item">Authors</router-link>
        <router-link v-if="!isAuthenticated" to="/login" class="navbar-item">Login</router-link>
        <router-link v-if="!isAuthenticated" to="/register" class="navbar-item">Register</router-link>
        <router-link v-if="isAuthenticated" to="/logout" class="navbar-item">Logout</router-link>
      </div>
    </nav>
  </div>
</template>

<style scoped>
.navbar {
  padding: 0.5rem;
  margin-bottom: 2rem;
}
</style>
