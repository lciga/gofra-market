import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'
import store from '../store'

const router = createRouter({
	history: createWebHistory(),
	routes,
})

router.beforeEach((to, from, next) => {
    const isAuthenticated = store.state.auth.isAuthenticated

    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login')
        return
    }

    if (to.meta.requiresGuest && isAuthenticated) {
        next('/')
        return
    }

    next()
})

export default router