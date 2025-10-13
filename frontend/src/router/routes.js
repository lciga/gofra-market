const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import('pages/Index.vue'),
    },
    {
        path: '/listing/:id',
        name: 'Listing Detail',
        component: () => import('pages/ListingDetail.vue'),
        props: true,
    },
    {
        path: '/profile',
        name: 'Profile',
        component: () => import('pages/Profile.vue'),
        meta: {requiresAuth: true},
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('pages/Login.vue'),
        meta: {requiresAuth: true},
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('pages/Register.vue'),
        meta: {requiresAuth: true},
    },
    {
        path: '/:catchAll(.*)*',
        component: () => import('pages/404.vue')
    },
]

export default routes