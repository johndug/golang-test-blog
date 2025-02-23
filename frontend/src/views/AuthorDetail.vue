<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthorsStore } from '../stores/authors'
import HeaderNavbar from '../components/HeaderNavbar.vue'

const route = useRoute()
const authorsStore = useAuthorsStore()
const loading = ref(false)
const error = ref('')

const fetchAuthor = async () => {
    try {
        loading.value = true
        const slug = route.params.slug as unknown as string
        await authorsStore.fetchAuthorBySlug(slug)
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Failed to load author'
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    fetchAuthor()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <HeaderNavbar />
    <main class="container mx-auto px-4 py-8">
      <h1 class="text-3xl font-bold mb-8">Author Detail</h1>
    </main>
    <div v-if="loading" class="text-center py-8 text-gray-500">Loading...</div>
    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
      {{ error }}
    </div>
    <div v-else-if="authorsStore.currentAuthor">
      <h2 class="text-2xl font-bold mb-4">{{ authorsStore.currentAuthor.first_name }} {{ authorsStore.currentAuthor.last_name }}</h2>
      <p class="text-gray-600 mb-4">{{ authorsStore.currentAuthor.bio }}</p>
    </div>
  </div>
</template>



