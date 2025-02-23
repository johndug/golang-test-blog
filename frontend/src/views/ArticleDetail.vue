<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useArticlesStore } from '../stores/articles'
import HeaderNavbar from '../components/HeaderNavbar.vue'

const route = useRoute()
const articlesStore = useArticlesStore()
const loading = ref(false)
const error = ref('')

const fetchArticle = async () => {
  try {
    loading.value = true
    const id = route.params.id as unknown as number
    await articlesStore.fetchArticleById(id)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load article'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchArticle()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <HeaderNavbar />
    <main class="container mx-auto px-4 py-8">
      <div v-if="loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500 mx-auto"></div>
      </div>

      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        {{ error }}
      </div>

      <article v-else-if="articlesStore.currentArticle" class="bg-white rounded-lg shadow-lg p-8">
        <h1 class="text-4xl font-bold mb-4">{{ articlesStore.currentArticle.title }}</h1>

        <div class="flex items-center text-gray-600 mb-8">
          <span>By {{ articlesStore.currentArticle.author?.first_name }} {{ articlesStore.currentArticle.author?.last_name }}</span>
          <span class="mx-2">•</span>
          <span>{{ new Date(articlesStore.currentArticle.created_at).toLocaleDateString() }}</span>
        </div>

        <div class="prose max-w-none">
          {{ articlesStore.currentArticle.content }}
        </div>
      </article>
    </main>
  </div>
</template>
