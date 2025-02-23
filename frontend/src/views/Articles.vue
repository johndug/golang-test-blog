<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useArticlesStore } from '../stores/articles'
import HeaderNavbar from '../components/HeaderNavbar.vue'

const articlesStore = useArticlesStore()
const loading = ref(false)
const error = ref('')

const fetchArticles = async () => {
  try {
    loading.value = true
    await articlesStore.fetchAllArticles()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load articles'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchArticles()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <HeaderNavbar />
    <main class="container mx-auto px-4 py-8">
      <h1 class="text-3xl font-bold mb-8">Articles</h1>

      <div v-if="loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500 mx-auto"></div>
      </div>

      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">T
        {{ error }}
      </div>

      <div v-else-if="articlesStore.articles.length === 0" class="text-center py-8 text-gray-500">
        No articles found.
      </div>

      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <article
          v-for="article in articlesStore.articles"
          :key="article.id"
          class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
        >
          <div class="p-6">
            <h2 class="text-xl font-semibold mb-2">
              <router-link
                :to="{ name: 'ArticleDetail', params: { id: article.id }}"
                class="text-indigo-600 hover:text-indigo-800"
              >
                {{ article.title }}
              </router-link>
            </h2>
            <p class="text-gray-600 mb-4">{{ article.short_description }}</p>
            <div class="flex items-center text-sm text-gray-500">
              <span>By {{ article.author?.first_name }} {{ article.author?.last_name }}</span>
              <span class="mx-2">•</span>
              <span>{{ new Date(article.created_at).toLocaleDateString() }}</span>
            </div>
          </div>
        </article>
      </div>
    </main>
  </div>
</template>


