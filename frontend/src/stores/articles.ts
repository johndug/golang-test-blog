import { defineStore } from 'pinia'
import axios from 'axios'
import type { Article } from '../types'

interface ArticlesState {
  articles: Article[]
  currentArticle: Article | null
  loading: boolean
  error: string | null
}

export const useArticlesStore = defineStore('articles', {
  state: (): ArticlesState => ({
    articles: [],
    currentArticle: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchAllArticles() {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get<Article[]>('/api/articles')
        this.articles = response.data
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Failed to fetch articles'
        throw this.error
      } finally {
        this.loading = false
      }
    },

    async fetchArticleById(id: number) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get<Article>(`/api/articles/${id}`)
        this.currentArticle = response.data
        return this.currentArticle
      } catch (error: any) {
        this.error = error.response?.data?.error || 'Failed to fetch article'
        throw this.error
      } finally {
        this.loading = false
      }
    }
  }
})
