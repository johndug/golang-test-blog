import { defineStore } from 'pinia'
import axios from 'axios'
import type { Author } from '../types'

interface AuthorsState {
    authors: Author[]
    currentAuthor: Author | null
    loading: boolean
    error: string | null
}

export const useAuthorsStore = defineStore('authors', {
    state: (): AuthorsState => ({
        authors: [],
        currentAuthor: null,
        loading: false,
        error: null
    }),

    actions: {
        async fetchAllAuthors() {
            this.loading = true
            this.error = null
            try {
                const response = await axios.get<Author[]>('/api/authors')
                this.authors = response.data
            } catch (error: any) {
                this.error = error.response?.data?.error || 'Failed to fetch authors'
                throw this.error
            } finally {
                this.loading = false
            }
        },

        async fetchAuthorBySlug(slug: string) {
            this.loading = true
            this.error = null
            try {
                const response = await axios.get<Author>(`/api/authors/${slug}`)
                this.currentAuthor = response.data
            } catch (error: any) {
                this.error = error.response?.data?.error || 'Failed to fetch author'
                throw this.error
            } finally {
                this.loading = false
            }
        }
    }
})

