<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import HeaderNavbar from '../../components/HeaderNavbar.vue'
import { computed, onMounted } from 'vue'

const router = useRouter()
const authStore = useAuthStore()

const isAuthor = computed(() => authStore.user?.role?.name === 'author')
const isAdmin = computed(() => authStore.user?.role?.name === 'admin')

onMounted(() => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <HeaderNavbar />
    <main class="container mx-auto px-4 py-8">
      <h1 class="text-3xl font-bold mb-8">Dashboard</h1>
      <!-- TODO: Add and edit articles if you're an author -->
      <div v-if="isAuthor">
        <h2 class="text-2xl font-bold mb-4">Your Articles</h2>
        <!-- TODO: Add a list of articles -->
      </div>
      <!-- TODO: Add and edit authors if you're an admin -->
      <div v-if="isAdmin">
        <h2 class="text-2xl font-bold mb-4">Authors</h2>
        <!-- TODO: Add a list of authors -->
      </div>
    </main>

  </div>
</template>
