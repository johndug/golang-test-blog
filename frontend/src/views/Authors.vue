<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthorsStore } from '../stores/authors'
import HeaderNavbar from '../components/HeaderNavbar.vue'

const authorsStore = useAuthorsStore()
const loading = ref(false)
const error = ref('')

const fetchAuthors = async () => {
  try {
    loading.value = true
    await authorsStore.fetchAllAuthors()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load authors'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAuthors()
})
</script>

<template>

    <div>
        <HeaderNavbar />
        <main class="container mx-auto px-4 py-8">
            <h1 class="text-3xl font-bold mb-8">Authors</h1>
            <ul>
                <li v-for="author in authorsStore.authors" :key="author.id">
                    <router-link :to="{ name: 'AuthorDetail', params: { slug: author.slug } }">
                        {{ author.first_name }} {{ author.last_name }}
                    </router-link>
                </li>
            </ul>
        </main>
    </div>
</template>
