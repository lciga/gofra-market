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
        path: '/admin',
        name: 'Admin',
        component: () => import('pages/Admin.vue'),
        meta: {requiresAuth: true, role: 'admin'},
    },
    {
        path: '/content-review',
        name: 'Content Review',
        component: () => import('pages/ContentReview.vue'),
        meta: {requiresAuth: true},
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('pages/Login.vue'),
        meta: {requiresGuest: true},
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('pages/Register.vue'),
        meta: {requiresGuest: true},
    },
    {
        path: '/:catchAll(.*)*',
        component: () => import('pages/Error404.vue')
    },
]

export default routes
