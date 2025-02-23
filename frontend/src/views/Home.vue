<script setup lang="ts">
import HeaderNavbar from '../components/HeaderNavbar.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <HeaderNavbar />
    <main class="container mx-auto px-4 py-8">
      <h1 class="text-2xl font-bold mb-4">Welcome to the Home Page</h1>
    </main>
  </div>
</template>

<style scoped>
.container {
  padding: 2rem;
}
</style>
