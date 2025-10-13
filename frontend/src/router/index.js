import {createRouter, createWebHistory} from 'vue-router'
import routers from './routes'
import store from '../store'

const router = createRouter({
    history: createWebHistory(),
    routers
})

router.beforeEach((to, from, next) => {
    const isAuthenticated = store.state.auth.isAuthenticated
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login')
    } else if (to.meta.requiresAuth && isAuthenticated) {
        next('/')
    } else {
        next()
    }
})

export default router