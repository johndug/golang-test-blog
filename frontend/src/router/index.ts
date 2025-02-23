import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Articles from '../views/Articles.vue'
import ArticleDetail from '../views/ArticleDetail.vue'
import Authors from '../views/Authors.vue'
import AuthorDetail from '../views/AuthorDetail.vue'
import Dashboard from '../views/Authenticated/dashboard.vue'
const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/articles',
      name: 'Articles',
      component: Articles
    },
    {
      path: '/authors',
      name: 'Authors',
      component: Authors
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { guest: true }
    },
    {
      path: '/register',
      name: 'Register',
      component: Register,
      meta: { guest: true }
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
      meta: { requiresAuth: true }
    },
    {
      path: '/logout',
      name: 'Logout',
      component: {
        beforeRouteEnter(to, from, next) {
          const authStore = useAuthStore()
          authStore.logout()
          next('/')
        }
      }
    },

    {
      path: '/articles/:id',
      name: 'ArticleDetail',
      component: ArticleDetail
    },
    {
      path: '/authors/:slug',
      name: 'AuthorDetail',
      component: AuthorDetail
    }
  ]
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }

  // Prevent authenticated users from accessing guest routes
  if (to.meta.guest && authStore.isAuthenticated) {
    next('/')
    return
  }

  next()
})

export default router
